package postgres

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"gobanking-v2/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool
}

func NewPool(cfg config.Config) (*Pool, error) {
	// pgx uses URL-encoded DSN.
	// Trade-off: using a URL DSN avoids manual escaping.
	q := url.Values{}
	q.Set("sslmode", cfg.DBSSLMode)

	// Example: postgres://user:pass@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?%s",
		url.PathEscape(cfg.DBUser),
		url.PathEscape(cfg.DBPassword),
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		q.Encode(),
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 10 * time.Minute

	p, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Ping(ctx); err != nil {
		p.Close()
		return nil, err
	}

	return &Pool{pool: p}, nil
}

func (p *Pool) Close() error {
	if p == nil || p.pool == nil {
		return nil
	}
	p.pool.Close()
	return nil

}

func (p *Pool) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

var _ = time.Second
