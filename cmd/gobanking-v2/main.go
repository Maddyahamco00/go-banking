package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gobanking-v2/config"
	"gobanking-v2/delivery/http"
	"gobanking-v2/infrastructure/postgres"
	"gobanking-v2/internal/app"
	healthuc "gobanking-v2/internal/health"
	"gobanking-v2/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	_ = config.LoadDotEnv() // best-effort for local dev

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	lg, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	sugar := lg
	sugar.Info("startup",
		zap.String("app_name", cfg.AppName),
		zap.String("app_env", cfg.AppEnv),
		zap.Int("app_port", cfg.AppPort),
		zap.String("log_level", cfg.LogLevel),
	)

	db, err := postgresInit(cfg)
	if err != nil {
		sugar.Fatal("db_init_failed", zap.Error(err))
	}

	healthHandler := healthuc.NewHandler(healthuc.New(cfg.AppName))
	router := http.NewRouter(healthHandler)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application := app.New(cfg, lg.Logger, db)
	if err := application.Run(ctx, router); err != nil {
		// ListenAndServe returns http.ErrServerClosed on graceful shutdown.
		if err != nil {
			sugar.Error("server_error", zap.Error(err), zap.Duration("at", time.Since(time.Now())))
		}
	}

	sugar.Info("exiting")
}

func postgresInit(cfg config.Config) (*postgres.Pool, error) {
	return postgres.NewPool(cfg)
}
