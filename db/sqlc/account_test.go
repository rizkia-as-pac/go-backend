package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	curr "github.com/tech_school/simple_bank/utils/currency"
	"github.com/tech_school/simple_bank/utils/random"
)

// setiap test function di golang namanya harus dimulai dengan "Test" dan mengandung *testing.T sebagai input
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	createdUser := createRandomUser(t)

	var arg CreateAccountParams = CreateAccountParams{
		Owner:    createdUser.Username,
		Balance:  random.RandomMoney(),
		Currency: random.RandomCurrency(),
	}
	account, err := testStore.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createdUser.Username, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestGetAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	accountFromDB, err := testStore.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)

	require.Equal(t, createdAccount.ID, accountFromDB.ID)
	require.Equal(t, createdAccount.Owner, accountFromDB.Owner)
	require.Equal(t, createdAccount.Balance, accountFromDB.Balance)
	require.Equal(t, createdAccount.Currency, accountFromDB.Currency)

	require.WithinDuration(t, createdAccount.CreatedAt, accountFromDB.CreatedAt, time.Second)
}

func TestGetAccountForUpdate(t *testing.T) {
	createdAccount := createRandomAccount(t)

	accountFromDB, err := testStore.GetAccountForUpdate(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, accountFromDB)

	require.Equal(t, createdAccount.ID, accountFromDB.ID)
	require.Equal(t, createdAccount.Owner, accountFromDB.Owner)
	require.Equal(t, createdAccount.Balance, accountFromDB.Balance)
	require.Equal(t, createdAccount.Currency, accountFromDB.Currency)

	require.WithinDuration(t, createdAccount.CreatedAt, accountFromDB.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	deletedAccountFromDB, err := testStore.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, deletedAccountFromDB)

	require.Equal(t, createdAccount.ID, deletedAccountFromDB.ID)
	require.Equal(t, createdAccount.Owner, deletedAccountFromDB.Owner)
	require.Equal(t, createdAccount.Balance, deletedAccountFromDB.Balance)
	require.Equal(t, createdAccount.Currency, deletedAccountFromDB.Currency)

	accountFromDB, err := testStore.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
	require.Empty(t, accountFromDB)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: random.RandomMoney(),
	}

	updatedAccountFromDB, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccountFromDB)

	require.Equal(t, createdAccount.ID, updatedAccountFromDB.ID)
	require.Equal(t, createdAccount.Owner, updatedAccountFromDB.Owner)
	require.Equal(t, arg.Balance, updatedAccountFromDB.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccountFromDB.Currency)
}

func TestListAccount(t *testing.T) {
	createdUser := createRandomUser(t)
	cur := []string{curr.USD, curr.JPY, curr.RUB}

	for i := 0; i < 3; i++ {
		var arg CreateAccountParams = CreateAccountParams{
			Owner:   createdUser.Username,
			Balance: random.RandomMoney(),
			// Currency: random.RandomCurrency(), // kita tidak menggunkan random currency karna kemungkinan random menghasilkan currency yang sama untuk satu username sangat tinggi
			Currency: cur[i],
		}

		_, err := testStore.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
	}

	// tes limit 2
	arg := ListAccountParams{
		Owner:  createdUser.Username,
		Limit:  2,
		Offset: 0,
	}

	listAccount, err := testStore.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listAccount, 2)

	for _, account := range listAccount {
		require.NotEmpty(t, account)
		// require.Equal(t, user.Username, account.Owner)
	}

	// tes limit 3
	arg = ListAccountParams{
		Owner:  createdUser.Username,
		Limit:  3,
		Offset: 0,
	}

	listAccount, err = testStore.ListAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, listAccount, 3)

	for _, account := range listAccount {
		require.NotEmpty(t, account)
		require.Equal(t, createdUser.Username, account.Owner)

	}
}
