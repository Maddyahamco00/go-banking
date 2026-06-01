package unit

import (
	"testing"

	"gobanking-v2/domain/ledger"
)

func TestPostTransaction_UnbalancedEntriesRejected(t *testing.T) {
	accounts := ledger.TestAccounts()

	req := ledger.PostTransactionRequest{
		Reference:   "txn-ref-1",
		Description: "unbalanced",
		Debits: []ledger.Posting{
			{AccountID: accounts.AssetCash, AmountCents: 100},
		},
		Credits: []ledger.Posting{
			{AccountID: accounts.RevenueFees, AmountCents: 90},
		},
	}

	_, err := ledger.PostTransaction(req, ledger.NowUTC())
	if err == nil {
		t.Fatalf("expected error for unbalanced transaction")
	}
}

func TestPostTransaction_ZeroOrNegativeAmountsRejected(t *testing.T) {
	accounts := ledger.TestAccounts()

	t.Run("zero rejected", func(t *testing.T) {
		req := ledger.PostTransactionRequest{
			Reference:   "txn-ref-2",
			Description: "zero",
			Debits: []ledger.Posting{
				{AccountID: accounts.AssetCash, AmountCents: 0},
			},
			Credits: []ledger.Posting{
				{AccountID: accounts.RevenueFees, AmountCents: 0},
			},
		}
		_, err := ledger.PostTransaction(req, ledger.NowUTC())
		if err == nil {
			t.Fatalf("expected error for zero amounts")
		}
	})

	t.Run("negative rejected", func(t *testing.T) {
		req := ledger.PostTransactionRequest{
			Reference:   "txn-ref-3",
			Description: "negative",
			Debits: []ledger.Posting{
				{AccountID: accounts.AssetCash, AmountCents: -1},
			},
			Credits: []ledger.Posting{
				{AccountID: accounts.RevenueFees, AmountCents: 1},
			},
		}
		_, err := ledger.PostTransaction(req, ledger.NowUTC())
		if err == nil {
			t.Fatalf("expected error for negative amounts")
		}
	})
}

func TestPostTransaction_MissingAccountsRejected(t *testing.T) {
	accounts := ledger.TestAccounts()

	req := ledger.PostTransactionRequest{
		Reference:   "txn-ref-4",
		Description: "missing account",
		Debits: []ledger.Posting{
			{AccountID: 999999, AmountCents: 100},
		},
		Credits: []ledger.Posting{
			{AccountID: accounts.RevenueFees, AmountCents: 100},
		},
	}

	_, err := ledger.PostTransaction(req, ledger.NowUTC())
	if err == nil {
		t.Fatalf("expected error for missing account")
	}
}

func TestPostTransaction_DuplicateReferencesRejected(t *testing.T) {
	accounts := ledger.TestAccounts()

	store := ledger.NewInMemoryReferenceStore()
	now := ledger.NowUTC()

	req1 := ledger.PostTransactionRequest{
		Reference:   "txn-ref-dup",
		Description: "first",
		Debits: []ledger.Posting{{AccountID: accounts.AssetCash, AmountCents: 50}},
		Credits: []ledger.Posting{{AccountID: accounts.RevenueFees, AmountCents: 50}},
	}
	_, err := ledger.PostTransactionWithStore(req1, now, store)
	if err != nil {
		t.Fatalf("unexpected error for first post: %v", err)
	}

	req2 := ledger.PostTransactionRequest{
		Reference:   "txn-ref-dup",
		Description: "second",
		Debits: []ledger.Posting{{AccountID: accounts.AssetCash, AmountCents: 10}},
		Credits: []ledger.Posting{{AccountID: accounts.RevenueFees, AmountCents: 10}},
	}
	_, err = ledger.PostTransactionWithStore(req2, now, store)
	if err == nil {
		t.Fatalf("expected error for duplicate reference")
	}
}

func TestPostTransaction_BalancedHappyPathCreatesEntries(t *testing.T) {
	accounts := ledger.TestAccounts()
	store := ledger.NewInMemoryReferenceStore()
	var txID ledger.TransactionID

	req := ledger.PostTransactionRequest{
		Reference:   "txn-ref-5",
		Description: "happy path",
		Debits: []ledger.Posting{
			{AccountID: accounts.AssetCash, AmountCents: 30},
			{AccountID: accounts.AssetCash, AmountCents: 20},
		},
		Credits: []ledger.Posting{
			{AccountID: accounts.RevenueFees, AmountCents: 50},
		},
	}

	posted, err := ledger.PostTransactionWithStore(req, ledger.NowUTC(), store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if posted.Reference != req.Reference {
		t.Fatalf("expected reference %q, got %q", req.Reference, posted.Reference)
	}
	if len(posted.Entries) != 3 {
		t.Fatalf("expected 3 ledger entries (2 debits + 1 credit), got %d", len(posted.Entries))
	}
	_ = txID
}

