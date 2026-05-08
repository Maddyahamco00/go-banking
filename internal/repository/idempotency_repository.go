package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
)

// IdempotencyRepository stores idempotency keys and their responses
type IdempotencyRepository struct {
	db *sqlx.DB
}

// NewIdempotencyRepository creates a new IdempotencyRepository
func NewIdempotencyRepository(db *sqlx.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

// IdempotencyRecord represents a stored idempotency entry
type IdempotencyRecord struct {
	Key       string          `json:"key" db:"key"`
	Response  json.RawMessage `json:"response" db:"response"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	ExpiresAt time.Time      `json:"expires_at" db:"expires_at"`
}

// Get retrieves an existing idempotency record
func (r *IdempotencyRepository) Get(ctx context.Context, key string) (*IdempotencyRecord, error) {
	var record IdempotencyRecord
	query := `SELECT key, response, created_at, expires_at FROM idempotency_keys WHERE key = $1 AND expires_at > NOW()`
	err := r.db.GetContext(ctx, &record, query, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &record, err
}

// Set stores a new idempotency record
func (r *IdempotencyRepository) Set(ctx context.Context, key string, response json.RawMessage, ttlHours int) error {
	query := `
		INSERT INTO idempotency_keys (key, response, created_at, expires_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (key) DO UPDATE SET response = EXCLUDED.response
	`
	_, err := r.db.ExecContext(ctx, query, key, response, time.Now(), time.Now().Add(time.Duration(ttlHours)*time.Hour))
	return err
}

// Delete removes an idempotency key
func (r *IdempotencyRepository) Delete(ctx context.Context, key string) error {
	query := `DELETE FROM idempotency_keys WHERE key = $1`
	_, err := r.db.ExecContext(ctx, query, key)
	return err
}

// Cleanup removes expired idempotency keys
func (r *IdempotencyRepository) Cleanup(ctx context.Context) (int64, error) {
	query := `DELETE FROM idempotency_keys WHERE expires_at < NOW()`
	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}