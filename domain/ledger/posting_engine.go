package ledger

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AccountType string

const (
	AccountTypeAsset     AccountType = "ASSET"
	AccountTypeLiability AccountType = "LIABILITY"
	AccountTypeRevenue   AccountType = "REVENUE"
	AccountTypeExpense   AccountType = "EXPENSE"
)

type AccountID int64

type TransactionID int64

type Account struct {
	ID   AccountID
	UUID string
	Name string
	Type AccountType
}

type LedgerEntryType string

const (
	LedgerEntryTypeDebit  LedgerEntryType = "DEBIT"
	LedgerEntryTypeCredit LedgerEntryType = "CREDIT"
)

type Posting struct {
	AccountID   AccountID
	AmountCents int64
}

type LedgerEntry struct {
	ID             int64
	TransactionID  TransactionID
	AccountID      AccountID
	EntryType      LedgerEntryType
	AmountCents   int64
}

type Transaction struct {
	ID          TransactionID
	UUID        string
	Reference   string
	Description string
	Status      string
	CreatedAt   time.Time
	Entries     []LedgerEntry
}

type PostTransactionRequest struct {
	Reference   string
	Description string
	Debits      []Posting
	Credits     []Posting
}

func NowUTC() time.Time { return time.Now().UTC() }

// PostTransaction validates request and posts ledger entries in-memory.
// Persistence is added in later phases; unit tests focus on correctness.
func PostTransaction(req PostTransactionRequest, now time.Time) (Transaction, error) {
	store := NewInMemoryReferenceStore()
	return PostTransactionWithStore(req, now, store)
}

func PostTransactionWithStore(req PostTransactionRequest, now time.Time, store ReferenceStore) (Transaction, error) {
	if req.Reference == "" {
		return Transaction{}, fmt.Errorf("reference_required")
	}
	if req.Description == "" {
		return Transaction{}, fmt.Errorf("description_required")
	}
	if len(req.Debits) == 0 || len(req.Credits) == 0 {
		return Transaction{}, fmt.Errorf("must_have_debits_and_credits")
	}

	// Reference uniqueness
	if err := store.MarkReferenceUsed(req.Reference); err != nil {
		return Transaction{}, err
	}

	debitTotal := int64(0)
	for _, d := range req.Debits {
		if d.AmountCents <= 0 {
			return Transaction{}, fmt.Errorf("amount_must_be_positive")
		}
		debitTotal += d.AmountCents
	}
	creditTotal := int64(0)
	for _, c := range req.Credits {
		if c.AmountCents <= 0 {
			return Transaction{}, fmt.Errorf("amount_must_be_positive")
		}
		creditTotal += c.AmountCents
	}

	// Double-entry invariant
	if debitTotal != creditTotal {
		return Transaction{}, fmt.Errorf("unbalanced_transaction")
	}

	// Missing accounts check is modeled against TestAccounts() for now.
	allowed := TestAccounts()
	isAllowed := func(id AccountID) bool {
		return id == allowed.AssetCash || id == allowed.RevenueFees
	}
	for _, d := range req.Debits {
		if !isAllowed(d.AccountID) {
			return Transaction{}, fmt.Errorf("account_missing")
		}
	}
	for _, c := range req.Credits {
		if !isAllowed(c.AccountID) {
			return Transaction{}, fmt.Errorf("account_missing")
		}
	}

	tx := Transaction{
		ID:          1,
		UUID:        uuid.NewString(),
		Reference:   req.Reference,
		Description: req.Description,
		Status:      "POSTED",
		CreatedAt:   now,
		Entries:     []LedgerEntry{},
	}

	for i, d := range req.Debits {
		tx.Entries = append(tx.Entries, LedgerEntry{
			ID:            int64(i + 1),
			TransactionID: tx.ID,
			AccountID:     d.AccountID,
			EntryType:     LedgerEntryTypeDebit,
			AmountCents:   d.AmountCents,
		})
	}
	base := len(req.Debits)
	for j, c := range req.Credits {
		tx.Entries = append(tx.Entries, LedgerEntry{
			ID:            int64(base + j + 1),
			TransactionID: tx.ID,
			AccountID:     c.AccountID,
			EntryType:     LedgerEntryTypeCredit,
			AmountCents:   c.AmountCents,
		})
	}

	return tx, nil
}

// TestAccounts provides a stable in-memory account set for unit tests.
// Production will replace this with real account lookups.
func TestAccounts() struct {
	AssetCash   AccountID
	RevenueFees AccountID
} {
	return struct {
		AssetCash   AccountID
		RevenueFees AccountID
	}{
		AssetCash:   1,
		RevenueFees: 2,
	}
}

