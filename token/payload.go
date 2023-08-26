package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// berbagai error yang dikembalikan oleh VerifyToken function
var ErrorExpiredToken = errors.New("token sudah expired")
var ErrorInvalidToken = errors.New("token is invalid")

// Payload berisi data payload yang terkandung didalam token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`  // use to identify token owner
	IssuedAt  time.Time `json:"issued_at"` // time when the token was created
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload membuat payload baru berdasarkan spesifik username dan durasi
func NewPayload(username string, duration time.Duration) (tokenPayload *Payload, err error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenId,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

// mengecek apakah payload dari token valid (is expired or not) atau tidak
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
