package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-banking/internal/domain"
)

// LedgerRepository handles ledger entry persistence
type LedgerRepository struct {
	db *sqlx.DB
}

// NewLedgerRepository creates a new LedgerRepository
func NewLedgerRepository(db *sqlx.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

// Create creates a new ledger entry (must be within a transaction)
func (r *LedgerRepository) Create(ctx context.Context, tx *sqlx.Tx, entry *domain.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (id, transaction_id, account_id, entry_type, amount, currency, balance_before, balance_after, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := tx.ExecContext(ctx, query,
		entry.ID, entry.TransactionID, entry.AccountID, entry.EntryType,
		entry.Amount, entry.Currency, entry.BalanceBefore, entry.BalanceAfter, entry.CreatedAt,
	)
	return err
}

// GetByTransactionID retrieves all ledger entries for a transaction
func (r *LedgerRepository) GetByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]domain.LedgerEntry, error) {
	var entries []domain.LedgerEntry
	query := `SELECT * FROM ledger_entries WHERE transaction_id = $1 ORDER BY created_at`
	err := r.db.SelectContext(ctx, &entries, query, transactionID)
	return entries, err
}

// GetByAccountID retrieves ledger entries for an account with pagination
func (r *LedgerRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.LedgerEntry, error) {
	var entries []domain.LedgerEntry
	query := `
		SELECT * FROM ledger_entries
		WHERE account_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	err := r.db.SelectContext(ctx, &entries, query, accountID, limit, offset)
	return entries, err
}

// GetAccountBalanceAt retrieves the account balance at a specific time
func (r *LedgerRepository) GetAccountBalanceAt(ctx context.Context, accountID uuid.UUID, at time.Time) (string, error) {
	var balance string
	query := `
		SELECT COALESCE(balance_after, '0')::text
		FROM ledger_entries
		WHERE account_id = $1 AND created_at <= $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	err := r.db.GetContext(ctx, &balance, query, accountID, at)
	if err != nil {
		return "0", nil
	}
	return balance, nil
}

// CountByAccountID returns total entries for an account
func (r *LedgerRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM ledger_entries WHERE account_id = $1`
	err := r.db.GetContext(ctx, &count, query, accountID)
	return count, err
}