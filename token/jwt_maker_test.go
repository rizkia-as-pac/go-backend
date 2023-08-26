package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
	"github.com/tech_school/simple_bank/utils/random"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32, "abcedfghijklmnopqrstuvwxyz"))
	require.NoError(t, err)

	username := random.RandomOwner()
	duration := time.Minute // one minute

	IssuedAt := time.Now()
	expired_at := IssuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expired_at, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(random.RandomString(32, "abcedfghijklmnopqrstuvwxyz"))
	require.NoError(t, err)

	username := random.RandomOwner()

	token, payload, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	// require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredToken.Error())
	require.Nil(t, payload)

}

// melakukan testing ketahanan token terhadap serangan yang sering terjadi yaitu none algorithm pada header
func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// mensimulasikan attacker membuat token
	username := random.RandomOwner()
	payload, err := NewPayload(username, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	// UnsafeAllowNoneSignatureType mensimulasikan attacker membuat header tanpa sign method, thus mereka tidak mengirimkan secret key pada token mereka
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	// mensimulasikan pengecekan token yang dilakukan server
	maker, err := NewJWTMaker(random.RandomString(32, "abcedfghijklmnopqrstuvwxyz"))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)
}
