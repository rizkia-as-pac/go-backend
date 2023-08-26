package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// meskipun jwt tidak mengharuskan seberapa panjang secret key, namun lebih baik kita membuatnya cukup panjang untuk alasan keamanan
const minSecretKeySize = 32 // minimun panjang 32 karakter

// JWTMaker is JSON web token maker implement token maker interface
// untuk case kali ini digunakan symetric key algorithm untuk melakukan sign pada token
type JWTMaker struct {
	secretKey string
}

// NewJWTMaker membuat a new JWTMaker
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("ukuran key tidak valid: setidaknya harus %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	// NewWithClaims menerima 2 argument : sign method, claims (payload)
	// dalam kasus ini digunakan sign method HS256
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	// jwtToken.Sihgnedstring mengenerate token
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	// return token, payload, err
	return token, payload, err
}

// VerifyToken melakukan verifikasi pada token. jika valid akan mengirimkan payload yang ada dalam body dari token tersebut
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// getsigning algorithm. karna tipe nya masih interface{} maka lalu konversi ke spesifik implementasi. kita gunakan SigningMethodHMAC karna SigningMethodHS256 merupakan instansiasi dari SigningMethodHMAC
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// untuk mengetahui tipe error. kita harus mengkonversi tipe errornya
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrorExpiredToken) {
			return nil, ErrorExpiredToken
		}
		return nil, ErrorInvalidToken
	}

	// untuk mendapatkan payload data dengan mengkonversi jwtToken.Claims into Payload object
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	}
	return payload, nil
}
