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

	"github.com/Tangyd893/TMS-Go/backend/internal/bootstrap"
	"github.com/Tangyd893/TMS-Go/backend/internal/config"
	"github.com/Tangyd893/TMS-Go/backend/internal/platform/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)

	server := bootstrap.NewHTTPServer(cfg, log)

	go func() {
		log.Info("starting http server", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server stopped unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutting down http server")
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown http server", slog.Any("error", err))
		os.Exit(1)
	}
}
