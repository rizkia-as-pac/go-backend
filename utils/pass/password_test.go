package pass

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/random"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := random.RandomPassword()

	// test create hash password
	hashedPassword1, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	// test check password
	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err) // no error mean password is correct

	// test check wrong password
	wrongPassword := random.RandomPassword()
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.Error(t, err)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// test check if two generated hash password from same base password are different
	hashedPassword2, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)
	require.NotEqual(t, hashedPassword1, hashedPassword2)
}
