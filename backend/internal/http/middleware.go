package httpapi

import (
	"net/http"
	"strings"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/auth"
)

func (r *Router) requireClaims(w http.ResponseWriter, req *http.Request) (auth.Claims, bool) {
	token := bearerToken(req)
	claims, err := r.authService.VerifyToken(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Valid bearer token is required.")
		return auth.Claims{}, false
	}

	return claims, true
}

func (r *Router) requireAdmin(w http.ResponseWriter, req *http.Request) (auth.Claims, bool) {
	claims, ok := r.requireClaims(w, req)
	if !ok {
		return auth.Claims{}, false
	}
	if claims.Role != auth.RoleAdmin {
		writeError(w, http.StatusForbidden, "forbidden", "Admin role is required.")
		return auth.Claims{}, false
	}

	return claims, true
}

func bearerToken(req *http.Request) string {
	header := strings.TrimSpace(req.Header.Get("Authorization"))
	if header == "" {
		return ""
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}

func withCORS(allowedOrigins []string, next http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	allowAll := false
	for _, origin := range allowedOrigins {
		if origin == "*" {
			allowAll = true
			continue
		}
		allowed[origin] = struct{}{}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin")
		if origin != "" {
			if allowAll {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if _, ok := allowed[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Max-Age", "600")
		}

		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, req)
	})
}
