package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go-banking/internal/domain"
	"go-banking/internal/repository"
)

// WalletService handles wallet operations
type WalletService struct {
	walletRepo  *repository.WalletRepository
	accountRepo *repository.AccountRepository
	ledgerSvc   *LedgerService
	txnRepo     *repository.TransactionRepository
}

// NewWalletService creates a new WalletService
func NewWalletService(
	walletRepo *repository.WalletRepository,
	accountRepo *repository.AccountRepository,
	ledgerSvc *LedgerService,
	txnRepo *repository.TransactionRepository,
) *WalletService {
	return &WalletService{
		walletRepo:  walletRepo,
		accountRepo: accountRepo,
		ledgerSvc:   ledgerSvc,
		txnRepo:     txnRepo,
	}
}

// CreateWalletRequest represents a wallet creation request
type CreateWalletRequest struct {
	OwnerID  uuid.UUID `json:"owner_id" binding:"required"`
	Currency string    `json:"currency" binding:"required"`
}

// Create creates a new wallet for a user
func (s *WalletService) Create(ctx context.Context, req *CreateWalletRequest) (*domain.WalletDetails, error) {
	// Default currency to NGN if not specified
	if req.Currency == "" {
		req.Currency = "NGN"
	}

	// Check if wallet already exists
	exists, err := s.walletRepo.Exists(ctx, req.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check wallet existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("wallet already exists for this owner")
	}

	// Create account first
	account := &domain.Account{
		ID:          uuid.New(),
		OwnerID:     req.OwnerID,
		AccountType: domain.AccountTypeWallet,
		Currency:    req.Currency,
		Balance:     decimal.Zero,
		Tier:        domain.Tier1,
		Status:      domain.AccountStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Create wallet
	wallet := &domain.Wallet{
		ID:        uuid.New(),
		OwnerID:   req.OwnerID,
		AccountID: account.ID,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.walletRepo.Create(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &domain.WalletDetails{
		Wallet:      *wallet,
		Balance:     account.Balance,
		Currency:    account.Currency,
		Tier:        account.Tier,
		Status:      account.Status,
	}, nil
}

// GetWallet retrieves wallet details by owner ID
func (s *WalletService) GetWallet(ctx context.Context, ownerID uuid.UUID) (*domain.WalletDetails, error) {
	details, err := s.walletRepo.GetDetails(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if details == nil {
		return nil, fmt.Errorf("wallet not found")
	}
	return details, nil
}

// GetWalletByID retrieves wallet by wallet ID
func (s *WalletService) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*domain.WalletDetails, error) {
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("wallet not found")
	}

	details, err := s.walletRepo.GetDetails(ctx, wallet.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet details: %w", err)
	}
	return details, nil
}

// GetTransactions retrieves transaction history for a wallet
func (s *WalletService) GetTransactions(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]domain.Transaction, int, error) {
	// Get wallet to find account ID
	details, err := s.walletRepo.GetDetails(ctx, ownerID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Get transactions
	txns, err := s.txnRepo.GetByAccountID(ctx, details.AccountID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Get total count
	total, err := s.txnRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return txns, total, nil
}

// GetLedgerEntries retrieves ledger entries for a wallet's account
func (s *WalletService) GetLedgerEntries(ctx context.Context, ownerID uuid.UUID, limit, offset int) ([]domain.LedgerEntry, int, error) {
	details, err := s.walletRepo.GetDetails(ctx, ownerID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get wallet: %w", err)
	}

	return s.ledgerSvc.GetAccountLedgerEntries(ctx, details.AccountID, limit, offset)
}

// FundRequest represents a wallet funding request
type FundRequest struct {
	OwnerID       uuid.UUID
	ToAccountID   uuid.UUID
	Amount        decimal.Decimal
	Description   string
	IdempotencyKey string
}

// Fund adds money to a wallet (from external source - simplified for MVP)
func (s *WalletService) Fund(ctx context.Context, req *FundRequest) (*domain.Transaction, error) {
	// For MVP, this would typically come from a payment processor
	// Here we simulate a deposit from external source
	// In production, this would integrate with a payment gateway

	// Get the system account for funding (in production, this would be the payment processor account)
	systemAccountID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	result, err := s.ledgerSvc.Transfer(ctx, &TransferRequest{
		FromAccountID:  systemAccountID,
		ToAccountID:    req.ToAccountID,
		Amount:         req.Amount,
		Currency:       "NGN",
		Description:    req.Description,
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, err
	}

	return result.Transaction, nil
}