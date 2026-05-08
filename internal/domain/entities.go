package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccountType represents the type of account
type AccountType string

const (
	AccountTypeWallet AccountType = "wallet"
	AccountTypeEscrow AccountType = "escrow"
	AccountTypeSystem AccountType = "system"
)

// AccountStatus represents the status of an account
type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusSuspended AccountStatus = "suspended"
	AccountStatusClosed   AccountStatus = "closed"
)

// Tier represents KYC tier level
type Tier string

const (
	Tier1 Tier = "tier1"
	Tier2 Tier = "tier2"
)

// TierLimits defines transaction limits per tier
var TierLimits = map[Tier]decimal.Decimal{
	Tier1: decimal.NewFromInt(50000),   // 50,000 NGN
	Tier2: decimal.NewFromInt(5000000), // 5,000,000 NGN
}

// Account is the core balance-holding entity
type Account struct {
	ID          uuid.UUID      `json:"id" db:"id"`
	OwnerID     uuid.UUID      `json:"owner_id" db:"owner_id"`
	AccountType AccountType    `json:"account_type" db:"account_type"`
	Currency    string         `json:"currency" db:"currency"`
	Balance     decimal.Decimal `json:"balance" db:"balance"`
	Tier        Tier           `json:"tier" db:"tier"`
	Status      AccountStatus  `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

// CanDeduct checks if the account has sufficient balance
func (a *Account) CanDeduct(amount decimal.Decimal) bool {
	return a.Balance.GreaterThanOrEqual(amount)
}

// LedgerEntryType represents debit or credit
type LedgerEntryType string

const (
	LedgerEntryTypeDebit  LedgerEntryType = "debit"
	LedgerEntryTypeCredit LedgerEntryType = "credit"
)

// LedgerEntry is an immutable journal entry
type LedgerEntry struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	TransactionID  uuid.UUID       `json:"transaction_id" db:"transaction_id"`
	AccountID      uuid.UUID       `json:"account_id" db:"account_id"`
	EntryType      LedgerEntryType `json:"entry_type" db:"entry_type"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Currency       string          `json:"currency" db:"currency"`
	BalanceBefore  decimal.Decimal `json:"balance_before" db:"balance_before"`
	BalanceAfter   decimal.Decimal `json:"balance_after" db:"balance_after"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// TransactionType represents the type of financial operation
type TransactionType string

const (
	TransactionTypeTransfer       TransactionType = "transfer"
	TransactionTypeEscrowHold     TransactionType = "escrow_hold"
	TransactionTypeEscrowRelease  TransactionType = "escrow_release"
	TransactionTypeLoanDisbursement TransactionType = "loan_disbursement"
	TransactionTypeLoanRepayment  TransactionType = "loan_repayment"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

// Transaction is a financial operation record
type Transaction struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	IdempotencyKey string           `json:"idempotency_key,omitempty" db:"idempotency_key"`
	TransactionRef string           `json:"transaction_ref" db:"transaction_ref"`
	Type           TransactionType  `json:"type" db:"type"`
	Status         TransactionStatus `json:"status" db:"status"`
	FromAccountID  *uuid.UUID       `json:"from_account_id,omitempty" db:"from_account_id"`
	ToAccountID    *uuid.UUID       `json:"to_account_id,omitempty" db:"to_account_id"`
	Amount         decimal.Decimal  `json:"amount" db:"amount"`
	Currency       string           `json:"currency" db:"currency"`
	Description    string           `json:"description" db:"description"`
	Metadata       map[string]any  `json:"metadata,omitempty" db:"metadata"`
	ErrorMessage   string           `json:"error_message,omitempty" db:"error_message"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

// Wallet is the user-facing wallet entity
type Wallet struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	OwnerID   uuid.UUID  `json:"owner_id" db:"owner_id"`
	AccountID uuid.UUID  `json:"account_id" db:"account_id"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// WalletDetails combines wallet with account info
type WalletDetails struct {
	Wallet
	Balance    decimal.Decimal `json:"balance"`
	Currency   string          `json:"currency"`
	Tier       Tier            `json:"tier"`
	Status     AccountStatus   `json:"account_status"`
}

// EscrowReleaseTrigger defines how escrow is released
type EscrowReleaseTrigger string

const (
	EscrowReleaseTriggerManual               EscrowReleaseTrigger = "manual"
	EscrowReleaseTriggerDeliveryConfirmation EscrowReleaseTrigger = "delivery_confirmation"
)

// EscrowStatus represents the status of an escrow hold
type EscrowStatus string

const (
	EscrowStatusHeld     EscrowStatus = "held"
	EscrowStatusReleased EscrowStatus = "released"
)

// EscrowHold tracks funds held in escrow
type EscrowHold struct {
	ID             uuid.UUID            `json:"id" db:"id"`
	TransactionID  uuid.UUID            `json:"transaction_id" db:"transaction_id"`
	FromAccountID  uuid.UUID            `json:"from_account_id" db:"from_account_id"`
	ToAccountID    uuid.UUID            `json:"to_account_id" db:"to_account_id"`
	Amount         decimal.Decimal      `json:"amount" db:"amount"`
	ReleaseTrigger EscrowReleaseTrigger `json:"release_trigger" db:"release_trigger"`
	Status         EscrowStatus         `json:"status" db:"status"`
	ReleasedAt     *time.Time          `json:"released_at,omitempty" db:"released_at"`
	CreatedAt      time.Time            `json:"created_at" db:"created_at"`
}

// LoanStatus represents the status of a loan
type LoanStatus string

const (
	LoanStatusPending    LoanStatus = "pending"
	LoanStatusDisbursed  LoanStatus = "disbursed"
	LoanStatusRepaying   LoanStatus = "repaying"
	LoanStatusDefaulted  LoanStatus = "defaulted"
	LoanStatusPaid       LoanStatus = "paid"
)

// Loan represents a micro-loan record
type Loan struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	OwnerID       uuid.UUID       `json:"owner_id" db:"owner_id"`
	Principal     decimal.Decimal `json:"principal" db:"principal"`
	InterestRate  decimal.Decimal `json:"interest_rate" db:"interest_rate"`
	TotalDue      decimal.Decimal `json:"total_due" db:"total_due"`
	AmountPaid    decimal.Decimal `json:"amount_paid" db:"amount_paid"`
	Status        LoanStatus      `json:"status" db:"status"`
	DueDate       time.Time       `json:"due_date" db:"due_date"`
	DisbursedAt   *time.Time     `json:"disbursed_at,omitempty" db:"disbursed_at"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

// KYCIDType represents the type of ID used for KYC
type KYCIDType string

const (
	KYCIDTypeBVN KYCIDType = "bvn"
	KYCIDTypeNIN KYCIDType = "nin"
)

// KYCStatus represents the KYC verification status
type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusVerified KYCStatus = "verified"
	KYCStatusRejected KYCStatus = "rejected"
)

// KYCRecord holds KYC verification data
type KYCRecord struct {
	ID              uuid.UUID   `json:"id" db:"id"`
	OwnerID         uuid.UUID   `json:"owner_id" db:"owner_id"`
	IDType          KYCIDType   `json:"id_type" db:"id_type"`
	IDNumber        string      `json:"id_number" db:"id_number"`
	Tier            Tier        `json:"tier" db:"tier"`
	VerificationRef string      `json:"verification_ref,omitempty" db:"verification_ref"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	Status          KYCStatus   `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}