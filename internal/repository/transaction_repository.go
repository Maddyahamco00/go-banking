package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-banking/internal/domain"
)

// TransactionRepository handles transaction persistence
type TransactionRepository struct {
	db *sqlx.DB
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create creates a new transaction (should be within a DB transaction)
func (r *TransactionRepository) Create(ctx context.Context, tx *sqlx.Tx, txn *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, idempotency_key, transaction_ref, type, status, from_account_id, to_account_id, amount, currency, description, metadata, created_at, updated_at)
		VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	var metadataJSON []byte
	if txn.Metadata != nil {
		metadataJSON, _ = json.Marshal(txn.Metadata)
	}

	_, err := tx.ExecContext(ctx, query,
		txn.ID, txn.IdempotencyKey, txn.TransactionRef, txn.Type, txn.Status,
		txn.FromAccountID, txn.ToAccountID, txn.Amount, txn.Currency,
		txn.Description, metadataJSON, txn.CreatedAt, txn.UpdatedAt,
	)
	return err
}

// UpdateStatus updates transaction status
func (r *TransactionRepository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, status domain.TransactionStatus, errorMsg string) error {
	query := `UPDATE transactions SET status = $1, error_message = $2, updated_at = $3 WHERE id = $4`
	_, err := tx.ExecContext(ctx, query, status, errorMsg, time.Now(), id)
	return err
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	var txn domain.Transaction
	query := `SELECT * FROM transactions WHERE id = $1`
	err := r.db.GetContext(ctx, &txn, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &txn, err
}

// GetByIdempotencyKey retrieves a transaction by idempotency key
func (r *TransactionRepository) GetByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	var txn domain.Transaction
	query := `SELECT * FROM transactions WHERE idempotency_key = $1`
	err := r.db.GetContext(ctx, &txn, query, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &txn, err
}

// GetByTransactionRef retrieves a transaction by reference
func (r *TransactionRepository) GetByTransactionRef(ctx context.Context, ref string) (*domain.Transaction, error) {
	var txn domain.Transaction
	query := `SELECT * FROM transactions WHERE transaction_ref = $1`
	err := r.db.GetContext(ctx, &txn, query, ref)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &txn, err
}

// GetByAccountID retrieves transactions involving an account with pagination
func (r *TransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.Transaction, error) {
	var txns []domain.Transaction
	query := `
		SELECT * FROM transactions
		WHERE from_account_id = $1 OR to_account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &txns, query, accountID, limit, offset)
	return txns, err
}

// GetAll retrieves all transactions with pagination (admin)
func (r *TransactionRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Transaction, error) {
	var txns []domain.Transaction
	query := `SELECT * FROM transactions ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &txns, query, limit, offset)
	return txns, err
}

// Count returns total transaction count
func (r *TransactionRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM transactions`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}

// GetByStatus retrieves transactions by status
func (r *TransactionRepository) GetByStatus(ctx context.Context, status domain.TransactionStatus, limit, offset int) ([]domain.Transaction, error) {
	var txns []domain.Transaction
	query := `SELECT * FROM transactions WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	err := r.db.SelectContext(ctx, &txns, query, status, limit, offset)
	return txns, err
}

// GetByOwnerID retrieves all transactions for a user's accounts
func (r *TransactionRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]domain.Transaction, error) {
	var txns []domain.Transaction
	query := `
		SELECT t.* FROM transactions t
		JOIN accounts a ON t.from_account_id = a.id OR t.to_account_id = a.id
		WHERE a.owner_id = $1
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &txns, query, ownerID, limit, offset)
	return txns, err
}