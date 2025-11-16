package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NetPo4ki/pull-review/internal/app"
	"github.com/NetPo4ki/pull-review/internal/config"
	"github.com/NetPo4ki/pull-review/internal/log"
)

func main() {
	cfg := config.Load()
	logger := log.NewLogger(cfg.AppEnv, cfg.LogLevel)

	handler := app.NewRouter(logger)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server_listen", "addr", srv.Addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server_error", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server_shutdown_error", "err", err)
	} else {
		logger.Info("server_shutdown_complete")
	}
}
