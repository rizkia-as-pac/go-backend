package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/pass"
	"github.com/tech_school/simple_bank/utils/random"
)

func createRandomUser(t *testing.T) User {

	person := random.RandomPerson()

	hashedPassword, err := pass.HashedPassword(person.Password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       person.Username,
		HashedPassword: hashedPassword,
		FullName:       person.FullName,
		Email:          person.Email,
	}

	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero()) // cek apakah field ini kosong apa tidak
	require.NotEmpty(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	userFromDB, err := testStore.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, userFromDB)

	require.Equal(t, createdUser.Username, userFromDB.Username)
	require.Equal(t, createdUser.HashedPassword, userFromDB.HashedPassword)
	require.Equal(t, createdUser.FullName, userFromDB.FullName)
	require.Equal(t, createdUser.Email, userFromDB.Email)

	require.WithinDuration(t, createdUser.CreatedAt, userFromDB.CreatedAt, time.Second) // cek dua waktu apakah terpisah jauh atau tidak
	require.WithinDuration(t, createdUser.CreatedAt, userFromDB.CreatedAt, time.Second) // cek dua waktu apakah terpisah jauh atau tidak
}
