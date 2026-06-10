package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/auth"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/config"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/database"
	httpapi "github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startupCancel()

	dbPool, err := database.Connect(startupCtx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := database.Migrate(startupCtx, dbPool); err != nil {
		logger.Error("failed to apply database migrations", "error", err)
		os.Exit(1)
	}

	authService := auth.NewService(auth.NewRepository(dbPool), cfg.JWTSecret, cfg.JWTTTL)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(cfg, logger, authService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("starting backend", "addr", server.Addr, "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	<-shutdownCtx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("backend stopped")
}
