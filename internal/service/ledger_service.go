package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"go-banking/internal/domain"
	"go-banking/internal/repository"
)

// LedgerService handles double-entry ledger operations
type LedgerService struct {
	db              *sqlx.DB
	accountRepo     *repository.AccountRepository
	ledgerRepo      *repository.LedgerRepository
	transactionRepo *repository.TransactionRepository
}

// NewLedgerService creates a new LedgerService
func NewLedgerService(
	db *sqlx.DB,
	accountRepo *repository.AccountRepository,
	ledgerRepo *repository.LedgerRepository,
	transactionRepo *repository.TransactionRepository,
) *LedgerService {
	return &LedgerService{
		db:              db,
		accountRepo:     accountRepo,
		ledgerRepo:      ledgerRepo,
		transactionRepo: transactionRepo,
	}
}

// TransferResult holds the result of a transfer operation
type TransferResult struct {
	Transaction  *domain.Transaction
	LedgerEntries []domain.LedgerEntry
}

// TransferRequest represents a transfer request
type TransferRequest struct {
	FromAccountID  uuid.UUID
	ToAccountID    uuid.UUID
	Amount         decimal.Decimal
	Currency       string
	Description    string
	IdempotencyKey string
	Metadata       map[string]any
}

// Transfer performs an atomic double-entry transfer between accounts
// This is the CRITICAL financial operation - all money movement goes through here
func (s *LedgerService) Transfer(ctx context.Context, req *TransferRequest) (*TransferResult, error) {
	// Generate transaction reference
	txnRef := fmt.Sprintf("TXN-%s", uuid.New().String()[:8])

	// Start database transaction
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create transaction record first (status = pending)
	txn := &domain.Transaction{
		ID:             uuid.New(),
		IdempotencyKey: req.IdempotencyKey,
		TransactionRef: txnRef,
		Type:           domain.TransactionTypeTransfer,
		Status:         domain.TransactionStatusPending,
		FromAccountID:  &req.FromAccountID,
		ToAccountID:    &req.ToAccountID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Description:    req.Description,
		Metadata:       req.Metadata,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err = s.transactionRepo.Create(ctx, tx, txn); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Lock and fetch source account (with row lock)
	fromAccount, err := s.accountRepo.GetForUpdate(ctx, tx, req.FromAccountID)
	if err != nil {
		return nil, fmt.Errorf("source account error: %w", err)
	}
	if fromAccount == nil {
		return nil, fmt.Errorf("source account not found")
	}
	if fromAccount.Status != domain.AccountStatusActive {
		return nil, fmt.Errorf("source account is not active: %s", fromAccount.Status)
	}

	// Lock and fetch destination account (with row lock)
	toAccount, err := s.accountRepo.GetForUpdate(ctx, tx, req.ToAccountID)
	if err != nil {
		return nil, fmt.Errorf("destination account error: %w", err)
	}
	if toAccount == nil {
		return nil, fmt.Errorf("destination account not found")
	}
	if toAccount.Status != domain.AccountStatusActive {
		return nil, fmt.Errorf("destination account is not active: %s", toAccount.Status)
	}

	// Validate sufficient balance
	if !fromAccount.CanDeduct(req.Amount) {
		updateErr := s.transactionRepo.UpdateStatus(ctx, tx, txn.ID, domain.TransactionStatusFailed, "insufficient funds")
		if updateErr != nil {
			return nil, fmt.Errorf("failed to update transaction status: %w", updateErr)
		}
		return nil, fmt.Errorf("insufficient funds in source account")
	}

	// Calculate new balances
	fromBalanceAfter := fromAccount.Balance.Sub(req.Amount)
	toBalanceAfter := toAccount.Balance.Add(req.Amount)

	// Create DEBIT ledger entry (source account, balance decreases)
	debitEntry := &domain.LedgerEntry{
		ID:            uuid.New(),
		TransactionID: txn.ID,
		AccountID:     fromAccount.ID,
		EntryType:     domain.LedgerEntryTypeDebit,
		Amount:        req.Amount,
		Currency:      req.Currency,
		BalanceBefore: fromAccount.Balance,
		BalanceAfter:  fromBalanceAfter,
		CreatedAt:     time.Now(),
	}
	if err = s.ledgerRepo.Create(ctx, tx, debitEntry); err != nil {
		return nil, fmt.Errorf("failed to create debit ledger entry: %w", err)
	}

	// Create CREDIT ledger entry (destination account, balance increases)
	creditEntry := &domain.LedgerEntry{
		ID:            uuid.New(),
		TransactionID: txn.ID,
		AccountID:     toAccount.ID,
		EntryType:     domain.LedgerEntryTypeCredit,
		Amount:        req.Amount,
		Currency:      req.Currency,
		BalanceBefore: toAccount.Balance,
		BalanceAfter:  toBalanceAfter,
		CreatedAt:     time.Now(),
	}
	if err = s.ledgerRepo.Create(ctx, tx, creditEntry); err != nil {
		return nil, fmt.Errorf("failed to create credit ledger entry: %w", err)
	}

	// Update source account balance
	if err = s.accountRepo.UpdateBalance(ctx, tx, fromAccount.ID, fromBalanceAfter); err != nil {
		return nil, fmt.Errorf("failed to update source account balance: %w", err)
	}

	// Update destination account balance
	if err = s.accountRepo.UpdateBalance(ctx, tx, toAccount.ID, toBalanceAfter); err != nil {
		return nil, fmt.Errorf("failed to update destination account balance: %w", err)
	}

	// Update transaction status to completed
	if err = s.transactionRepo.UpdateStatus(ctx, tx, txn.ID, domain.TransactionStatusCompleted, ""); err != nil {
		return nil, fmt.Errorf("failed to update transaction status: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &TransferResult{
		Transaction:   txn,
		LedgerEntries: []domain.LedgerEntry{*debitEntry, *creditEntry},
	}, nil
}

// GetLedgerEntries retrieves ledger entries for a transaction
func (s *LedgerService) GetLedgerEntries(ctx context.Context, transactionID uuid.UUID) ([]domain.LedgerEntry, error) {
	entries, err := s.ledgerRepo.GetByTransactionID(ctx, transactionID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// GetAccountLedgerEntries retrieves ledger entries for an account with pagination
func (s *LedgerService) GetAccountLedgerEntries(ctx context.Context, accountID uuid.UUID, limit, offset int) ([]domain.LedgerEntry, int, error) {
	entries, err := s.ledgerRepo.GetByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.ledgerRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return nil, 0, err
	}
	return entries, total, nil
}

// ValidateBalance validates that account balance matches ledger entries (audit function)
func (s *LedgerService) ValidateBalance(ctx context.Context, accountID uuid.UUID) (bool, decimal.Decimal, decimal.Decimal, error) {
	account, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil || account == nil {
		return false, decimal.Zero, decimal.Zero, err
	}

	// Sum all credits and debits for this account
	var totalDebits, totalCredits decimal.Decimal

	entries, err := s.ledgerRepo.GetByAccountID(ctx, accountID, 10000, 0)
	if err != nil {
		return false, decimal.Zero, decimal.Zero, err
	}

	for _, entry := range entries {
		if entry.EntryType == domain.LedgerEntryTypeDebit {
			totalDebits = totalDebits.Add(entry.Amount)
		} else {
			totalCredits = totalCredits.Add(entry.Amount)
		}
	}

	// Balance should equal credits - debits (for wallet accounts)
	expectedBalance := totalCredits.Sub(totalDebits)
	return account.Balance.Equal(expectedBalance), account.Balance, expectedBalance, nil
}