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

	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/catalog"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/catalogbot"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/config"
	"github.com/BratishkaDurovaTg/SWP-AromaType/backend/internal/database"
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

	bot := catalogbot.New(catalogbot.Config{
		Token:     os.Getenv("CATALOG_BOT_TOKEN"),
		Password:  os.Getenv("CATALOG_BOT_PASSWORD"),
		UploadDir: cfg.UploadDir,
	}, logger, catalog.NewRepository(dbPool))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := bot.Run(ctx); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("catalog bot stopped unexpectedly", "error", err)
		os.Exit(1)
	}
}
