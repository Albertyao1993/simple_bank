package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Albertyao1993/simple_bank/util"

	"github.com/stretchr/testify/require"
)

func creatRandomAccount(t *testing.T) Account {

	arg := CreateAccountParams{
		// Owner:    "Tom", // randomly generate ?
		// Balance:  1000,
		// Currency: "USD",
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	creatRandomAccount(t)
}

// test getaccount
func TestGetAccount(t *testing.T) {
	// create a few account
	account1 := creatRandomAccount(t)
	account2, err := testQueries.GetAccout(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {

	account1 := creatRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	testQueries.UpdateAccount(context.Background(), arg)

	err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)

	account2, err := testQueries.GetAccout(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, arg.Balance, account2.Balance)

	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := creatRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccout(context.Background(), account1.ID)
	require.Empty(t, account2)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())

}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		creatRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
