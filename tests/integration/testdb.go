package integration

import (
	"context"
	"fmt"
	"os"
	"time"

	"gobanking-v2/infrastructure/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

func loadDBConfigFromEnv() (DBConfig, error) {
	get := func(key string) (string, bool) {
		v, ok := os.LookupEnv(key)
		return v, ok
	}

	require := func(key string) (string, error) {
		v, ok := get(key)
		if !ok || v == "" {
			return "", fmt.Errorf("missing required env var %s", key)
		}
		return v, nil
	}

	host, err := require("DB_HOST")
	if err != nil {
		return DBConfig{}, err
	}
	portStr, err := require("DB_PORT")
	if err != nil {
		return DBConfig{}, err
	}
	name, err := require("DB_NAME")
	if err != nil {
		return DBConfig{}, err
	}
	user, err := require("DB_USER")
	if err != nil {
		return DBConfig{}, err
	}
	pass, err := require("DB_PASSWORD")
	if err != nil {
		return DBConfig{}, err
	}
	sslMode, ok := get("DB_SSLMODE")
	if !ok || sslMode == "" {
		sslMode = "disable"
	}

	var port int
	_, scanErr := fmt.Sscanf(portStr, "%d", &port)
	if scanErr != nil {
		return DBConfig{}, fmt.Errorf("DB_PORT must be an integer, got %q", portStr)
	}

	return DBConfig{
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: pass,
		SSLMode:  sslMode,
	}, nil
}

// OpenTestDB opens the application pgx pool.
// Cleanup/migrations are expected to be done by test bootstrap code.
func OpenTestDB() (*postgres.Pool, error) {
	cfg, err := loadDBConfigFromEnv()
	if err != nil {
		return nil, err
	}

	appCfg := struct {
		Host     string
		Port     int
		Name     string
		User     string
		Password string
		SSLMode  string
	}{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Name:     cfg.Name,
		User:     cfg.User,
		Password: cfg.Password,
		SSLMode:  cfg.SSLMode,
	}

	// Reuse existing postgres.NewPool via config.Config.
	// We keep it simple: only DB fields must be set for pool creation.
	// Other fields are not required by postgres.NewPool.
	//
	// Note: config.Config type lives in config package; we avoid importing here
	// to keep test helper focused. Instead, construct postgres pool directly.
	// However, postgres.NewPool requires config.Config, so we import it below.
	//
	// This function intentionally returns a helpful error if env is missing.

	// Inline import usage to avoid unused imports while keeping code straightforward.
	// (Go doesn't allow conditional imports, but we can keep imports minimal.)
	_ = appCfg

	// Proper implementation:
	// We'll create postgres.Pool using pgxpool directly to avoid dependence on config.Config.
	// But since postgres.NewPool already exists, we should use it.
	//
	// To avoid refactoring app code, we call pgxpool directly.
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		urlPathEscape(cfg.User),
		urlPathEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
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

	// postgres.Pool wraps pgxpool.Pool, but its internal field is unexported.
	// For now, return nil wrapper and rely on pgxpool directly from cleanup.
	// This helper is scaffolded; once cleanup/migrations are integrated properly,
	// we can refactor postgres.Pool to expose an accessor.
	return nil, fmt.Errorf("OpenTestDB not implemented: cannot construct postgres.Pool (unexported field)")
}

// CleanupAll truncates all user tables in public schema.
// This is a lightweight placeholder; for real schemas, prefer per-table cleanup.
func CleanupAll(ctx context.Context, pool *pgxpool.Pool) error {
	// Truncate in dependency-safe order.
	_, err := pool.Exec(ctx, `
		DO $$
		DECLARE r RECORD;
		BEGIN
			FOR r IN (SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname='public') LOOP
				-- no-op; just to force existence check
			END LOOP;
		END $$;
	`)
	// If schema has no tables, ignore cleanup.
	_ = err

	_, err = pool.Exec(ctx, `
		SELECT
			x.tablestr
		FROM (
			SELECT tablename || '"' AS tablestr
			FROM pg_catalog.pg_tables
			WHERE schemaname='public'
		) x;
	`)
	_ = err

	_, err = pool.Exec(ctx, `
		DO $$
		DECLARE stmt text;
		BEGIN
			SELECT string_agg(format('TRUNCATE TABLE public.%I RESTART IDENTITY CASCADE', tablename), '; ')
			INTO stmt
			FROM pg_catalog.pg_tables
			WHERE schemaname='public';

			IF stmt IS NULL THEN
				RETURN;
			END IF;

			execute stmt;
		END $$;
	`)
	return err
}

// urlPathEscape escapes minimal for DSN safety.
func urlPathEscape(s string) string {
	// Avoid pulling net/url; tests use simple env values (no spaces).
	// If values may contain special characters, this should be replaced by url.PathEscape.
	replacer := []struct{ old, new string }{
		{"@", "%40"},
		{":", "%3A"},
		{"/", "%2F"},
		{"?", "%3F"},
		{"#", "%23"},
	}
	out := s
	for _, r := range replacer {
		out = replaceAll(out, r.old, r.new)
	}
	return out
}

func replaceAll(s, old, new string) string {
	for {
		idx := indexOf(s, old)
		if idx < 0 {
			return s
		}
		s = s[:idx] + new + s[idx+len(old):]
	}
}

func indexOf(s, sub string) int {
	// naive search (small strings)
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
