package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/random"
)

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    random.RandomMoney(),
	}

	entry, err := testStore.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)

	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return entry
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	entryFromDB, err := testStore.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entryFromDB)

	require.Equal(t, entry.ID, entryFromDB.ID)
	require.Equal(t, entry.AccountID, entryFromDB.AccountID)
	require.Equal(t, entry.Amount, entryFromDB.Amount)

	require.WithinDuration(t, entry.CreatedAt, entryFromDB.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t)
	numEntries := 5

	for i := 0; i < numEntries; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    random.RandomMoney(),
		}

		_, err := testStore.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     int32(numEntries),
		Offset:    0,
	}

	entries, err := testStore.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, numEntries)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}
