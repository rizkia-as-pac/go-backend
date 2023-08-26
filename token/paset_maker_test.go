package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/random"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32, "abcedfghijklmnopqrstuvwxyz"))
	require.NoError(t, err)

	username := random.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	encryptedToken, payload,err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, encryptedToken)
	require.NotEmpty(t, payload)

	tokenPayload, err := maker.VerifyToken(encryptedToken)
	require.NoError(t, err)
	require.NotEmpty(t, tokenPayload)

	require.NotZero(t, tokenPayload.ID)
	require.Equal(t, username, tokenPayload.Username)
	require.WithinDuration(t, issuedAt, tokenPayload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, tokenPayload.ExpiredAt, time.Second)

}

func TestExpiredToken(t *testing.T) {
	maker, err := NewPasetoMaker(random.RandomString(32, "abcedfghijklmnopqrstuvwxyz"))
	require.NoError(t, err)

	username := random.RandomOwner()
	duration := time.Minute

	encryptedToken, payload, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, encryptedToken)
	require.NotEmpty(t, payload)

	tokenPayload, err := maker.VerifyToken(encryptedToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, tokenPayload)

}

// KITA TIDAK MEMBUTUHKAN TEST NONE ALGORITHM KARENA KASUS SEPERTI ITU TIDAK ADA DI PASETO
// TUGAS : MENULIS TEST INVALID TOKEN CASE
