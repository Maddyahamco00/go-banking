package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go-banking/internal/domain"
)

// WalletRepository handles wallet persistence
type WalletRepository struct {
	db *sqlx.DB
}

// NewWalletRepository creates a new WalletRepository
func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Create creates a new wallet
func (r *WalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	query := `
		INSERT INTO wallets (id, owner_id, account_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		wallet.ID, wallet.OwnerID, wallet.AccountID, wallet.Status, wallet.CreatedAt,
	)
	return err
}

// GetByID retrieves a wallet by ID
func (r *WalletRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `SELECT * FROM wallets WHERE id = $1`
	err := r.db.GetContext(ctx, &wallet, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &wallet, err
}

// GetByOwnerID retrieves a wallet by owner ID
func (r *WalletRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `SELECT * FROM wallets WHERE owner_id = $1`
	err := r.db.GetContext(ctx, &wallet, query, ownerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &wallet, err
}

// GetDetails retrieves wallet with full account details
func (r *WalletRepository) GetDetails(ctx context.Context, ownerID uuid.UUID) (*domain.WalletDetails, error) {
	var details domain.WalletDetails
	query := `
		SELECT
			w.id, w.owner_id, w.account_id, w.status, w.created_at,
			a.balance, a.currency, a.tier, a.status as account_status
		FROM wallets w
		JOIN accounts a ON w.account_id = a.id
		WHERE w.owner_id = $1
	`
	err := r.db.GetContext(ctx, &details, query, ownerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &details, err
}

// UpdateStatus updates wallet status
func (r *WalletRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE wallets SET status = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// Exists checks if a wallet exists for an owner
func (r *WalletRepository) Exists(ctx context.Context, ownerID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM wallets WHERE owner_id = $1)`
	err := r.db.GetContext(ctx, &exists, query, ownerID)
	return exists, err
}

// GetByAccountID retrieves a wallet by account ID
func (r *WalletRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `SELECT * FROM wallets WHERE account_id = $1`
	err := r.db.GetContext(ctx, &wallet, query, accountID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &wallet, err
}

// GetByOwnerIDTx retrieves a wallet by owner ID within a transaction
func (r *WalletRepository) GetByOwnerIDTx(ctx context.Context, tx *sqlx.Tx, ownerID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	query := `SELECT * FROM wallets WHERE owner_id = $1 FOR UPDATE`
	err := tx.GetContext(ctx, &wallet, query, ownerID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &wallet, err
}

// CountByAccountID returns the number of transactions for an account
func (r *WalletRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM transactions WHERE from_account_id = $1 OR to_account_id = $1`
	err := r.db.GetContext(ctx, &count, query, accountID)
	return count, err
}