package token

import "time"

type Maker interface {
	// CreateToken create new token untuk spesifik username dan durasi
	CreateToken(username string, duration time.Duration) (signedToken string, payload *Payload, err error)

	// VerifyToken melakukan verifikasi pada token. jika valid akan mengirimkan payload yang ada dalam body dari token tersebut
	VerifyToken(token string) (tokenPayload *Payload, err error)
}
