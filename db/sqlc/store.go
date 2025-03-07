package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store 接口定义
type Store interface {
	Querier // ✅ 正确的嵌入方式
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to excute db queries and transactions.
type SQLStore struct {
	*Queries
	db *sql.DB
}

// New Store creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

var txKey = struct{}{}

// exceTx executes a function within a database transanction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rollBackErr := tx.Rollback(); rollBackErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollBackErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferParams contains the input parameters of the transfer transanction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult contains the result of the transfer transanction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	if arg.FromAccountID == arg.ToAccountID {
		return result, fmt.Errorf("cannot transfer to same account")
	}

	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		// txName := ctx.Value(txKey)
		// fmt.Println(txName, "create transfer")

		// create a transfer record
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		//add account entries  FromEntry and ToEntry
		// fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		//add account entries
		// fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO : update accounts' balance  ----- it cover lock operaration

		// 获取原始账户信息
		// fmt.Println(txName, "get account 1")
		account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		// 更新账户
		// fmt.Println(txName, "update account 1 balance")
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: account1.Balance - arg.Amount,
		})
		if err != nil {
			return err
		}

		// 验证更新后的账户信息
		result.FromAccount, err = q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		// 可以添加余额验证
		if result.FromAccount.Balance != account1.Balance-arg.Amount {
			return fmt.Errorf("balance not updated correctly")
		}

		// 获取第二个账户信息
		// fmt.Println(txName, "get account 2")
		account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		// 更新第二个账户
		// fmt.Println(txName, "update account 2 balance")
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: account2.Balance + arg.Amount,
		})
		if err != nil {
			return err
		}

		// 验证更新后的账户信息
		result.ToAccount, err = q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		// 可以添加余额验证
		if result.ToAccount.Balance != account2.Balance+arg.Amount {
			return fmt.Errorf("balance not updated correctly")
		}

		return nil
	})
	return result, err

}
