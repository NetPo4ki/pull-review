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
	"github.com/NetPo4ki/pull-review/internal/repo"
	prssvc "github.com/NetPo4ki/pull-review/internal/service/prs"
	teamssvc "github.com/NetPo4ki/pull-review/internal/service/teams"
	userssvc "github.com/NetPo4ki/pull-review/internal/service/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	logger := log.NewLogger(cfg.AppEnv, cfg.LogLevel)

	pool, err := pgxpool.New(context.Background(), cfg.DBDSN)
	if err != nil {
		logger.Error("db_pool_new_error", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	store := repo.NewStore(pool)
	teamsService := teamssvc.New(store, store)
	usersService := userssvc.New(store, store)
	tx := repo.NewTxManager(pool)
	prService := prssvc.New(store, store, tx)

	handler := app.NewRouter(logger, teamsService, usersService, prService)

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
