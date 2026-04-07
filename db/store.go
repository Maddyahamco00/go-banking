package db

import (
	"context"
	"database/sql"
	"fmt"

	sqlc "github.com/Maddyahamco00/go-banking/db/sqlc"
)

type Store interface {
	sqlc.Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	DepositTx(ctx context.Context, arg DepositTxParams) (DepositTxResult, error)
	WithdrawTx(ctx context.Context, arg WithdrawTxParams) (WithdrawTxResult, error)
}

type SQLStore struct {
	*sqlc.Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := sqlc.New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

//
// ---------- TRANSFER ----------
//

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    sqlc.Transfer `json:"transfer"`
	FromAccount sqlc.Account  `json:"from_account"`
	ToAccount   sqlc.Account  `json:"to_account"`
	FromEntry   sqlc.Entry    `json:"from_entry"`
	ToEntry     sqlc.Entry    `json:"to_entry"`
}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, sqlc.CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(
				ctx, q,
				arg.FromAccountID, -arg.Amount,
				arg.ToAccountID, arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(
				ctx, q,
				arg.ToAccountID, arg.Amount,
				arg.FromAccountID, -arg.Amount,
			)
		}

		return err
	})

	return result, err
}

//
// ---------- DEPOSIT ----------
//

type DepositTxParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type DepositTxResult struct {
	Account sqlc.Account `json:"account"`
	Entry   sqlc.Entry   `json:"entry"`
}

func (store *SQLStore) DepositTx(ctx context.Context, arg DepositTxParams) (DepositTxResult, error) {
	var result DepositTxResult

	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		result.Entry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		result.Account, err = q.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
			ID:     arg.AccountID,
			Balance: arg.Amount,
		})

		return err
	})

	return result, err
}

//
// ---------- WITHDRAW ----------
//

type WithdrawTxParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

type WithdrawTxResult struct {
	Account sqlc.Account `json:"account"`
	Entry   sqlc.Entry   `json:"entry"`
}

func (store *SQLStore) WithdrawTx(ctx context.Context, arg WithdrawTxParams) (WithdrawTxResult, error) {
	var result WithdrawTxResult

	err := store.execTx(ctx, func(q *sqlc.Queries) error {
		var err error

		result.Entry, err = q.CreateEntry(ctx, sqlc.CreateEntryParams{
			AccountID: arg.AccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.Account, err = q.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
			ID:     arg.AccountID,
			Balance: -arg.Amount,
		})

		return err
	})

	return result, err
}

//
// ---------- HELPER ----------
//

func addMoney(
	ctx context.Context,
	q *sqlc.Queries,
	accountID1 int64, amount1 int64,
	accountID2 int64, amount2 int64,
) (account1 sqlc.Account, account2 sqlc.Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		ID:     accountID1,
		Balance: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, sqlc.UpdateAccountBalanceParams{
		ID:     accountID2,
		Balance: amount2,
	})

	return
}