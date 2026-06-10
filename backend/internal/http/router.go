package httpapi

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/config"
)

const serviceVersion = "0.1.0"

type Router struct {
	cfg    config.Config
	logger *slog.Logger
}

func NewRouter(cfg config.Config, logger *slog.Logger) http.Handler {
	router := &Router{
		cfg:    cfg,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", router.handleHealth)
	mux.HandleFunc("GET /docs", router.handleDocs)
	mux.HandleFunc("GET /openapi.yaml", router.handleOpenAPI)

	return logRequests(logger, mux)
}

func (r *Router) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "aromatype-backend",
		"version": serviceVersion,
	})
}

func (r *Router) handleOpenAPI(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, r.cfg.OpenAPIPath)
}

func (r *Router) handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := docsTemplate.Execute(w, nil); err != nil {
		r.logger.Error("failed to render docs page", "error", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		logger.Info("request completed", "method", req.Method, "path", req.URL.Path)
	})
}

var docsTemplate = template.Must(template.New("docs").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>AromaType API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/openapi.yaml",
        dom_id: "#swagger-ui"
      });
    </script>
  </body>
</html>`))
