package app

import (
	"context"
	"net/http"
	"time"

	"gobanking-v2/config"
	"gobanking-v2/infrastructure/postgres"
	"gobanking-v2/pkg/logger"

	"go.uber.org/zap"
)

type App struct {
	cfg     config.Config
	lg      *zap.Logger
	db      *postgres.Pool
	httpSrv *http.Server
}

func New(cfg config.Config, lg *zap.Logger, db *postgres.Pool) *App {
	return &App{cfg: cfg, lg: lg, db: db}
}

func (a *App) Run(ctx context.Context, handler http.Handler) error {
	addr := ":" + itoa(a.cfg.AppPort)

	a.httpSrv = &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	a.lg.Info("server_start", zap.String("addr", addr), zap.String("env", a.cfg.AppEnv))

	// Serve in goroutine.
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		a.lg.Info("shutdown_triggered")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = a.httpSrv.Shutdown(shutdownCtx)
		_ = a.db.Close()
		a.lg.Info("shutdown_complete")
		return nil
	case err := <-errCh:
		return err
	}
}

func itoa(i int) string {
	// tiny helper to avoid fmt import in core wiring
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var b [32]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		b[pos] = '-'
	}
	return string(b[pos:])
}

// Stop is handled via context cancellation in Run.
func (a *App) Stop() {}

// NewLoggerStart/stop logs are in main.
var _ = logger.New
var _ = time.Second
