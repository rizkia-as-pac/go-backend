package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// pasetomaker is a PASETO token maker that implement token maker interface
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker create a new PasetoMaker instance
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// paseto version 2 use chacha poly algorithm to encrypt the payload
	// chacha20poly1305.KeySize ukurannya adalah 32 karakter
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Ukuran key tidak benar : ukuran key harus tepat %d karakter", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (encryptedToken string, payload *Payload, err error) {
	tokenPayload, err := NewPayload(username, duration)
	if err != nil {
		return "", tokenPayload, err
	}

	encryptedToken, err = maker.paseto.Encrypt(maker.symmetricKey, tokenPayload, nil)
	return encryptedToken, tokenPayload, err
}

// VerifyToken melakukan verifikasi pada token. jika valid akan mengirimkan payload yang ada dalam body dari token tersebut
func (maker *PasetoMaker) VerifyToken(token string) (tokenPayload *Payload, err error) {
	payload := &Payload{} // create empty payload object

	err = maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrorInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
