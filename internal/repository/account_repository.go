package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"go-banking/internal/domain"
)

// AccountRepository handles account persistence
type AccountRepository struct {
	db *sqlx.DB
}

// NewAccountRepository creates a new AccountRepository
func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Create creates a new account
func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (id, owner_id, account_type, currency, balance, tier, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		account.ID, account.OwnerID, account.AccountType, account.Currency,
		account.Balance, account.Tier, account.Status, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT * FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &account, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}

// GetByOwnerID retrieves accounts by owner ID
func (r *AccountRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `SELECT * FROM accounts WHERE owner_id = $1`
	err := r.db.SelectContext(ctx, &accounts, query, ownerID)
	return accounts, err
}

// GetByOwnerIDAndType retrieves an account by owner ID and type
func (r *AccountRepository) GetByOwnerIDAndType(ctx context.Context, ownerID uuid.UUID, accountType domain.AccountType) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT * FROM accounts WHERE owner_id = $1 AND account_type = $2`
	err := r.db.GetContext(ctx, &account, query, ownerID, accountType)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}

// UpdateBalance updates account balance (must be within a transaction)
func (r *AccountRepository) UpdateBalance(ctx context.Context, tx *sqlx.Tx, id uuid.UUID, newBalance decimal.Decimal) error {
	query := `UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err := tx.ExecContext(ctx, query, newBalance, time.Now(), id)
	return err
}

// UpdateTier updates account tier
func (r *AccountRepository) UpdateTier(ctx context.Context, id uuid.UUID, tier domain.Tier) error {
	query := `UPDATE accounts SET tier = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, tier, time.Now(), id)
	return err
}

// UpdateStatus updates account status
func (r *AccountRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AccountStatus) error {
	query := `UPDATE accounts SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

// GetForUpdate retrieves account with row lock for update within transaction
func (r *AccountRepository) GetForUpdate(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account
	query := `SELECT * FROM accounts WHERE id = $1 FOR UPDATE`
	err := tx.GetContext(ctx, &account, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account not found: %s", id)
	}
	return &account, err
}

// GetAll retrieves all accounts (admin only)
func (r *AccountRepository) GetAll(ctx context.Context, limit, offset int) ([]domain.Account, error) {
	var accounts []domain.Account
	query := `SELECT * FROM accounts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &accounts, query, limit, offset)
	return accounts, err
}

// Count returns total number of accounts
func (r *AccountRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM accounts`
	err := r.db.GetContext(ctx, &count, query)
	return count, err
}