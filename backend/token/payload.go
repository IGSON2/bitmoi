package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrExpiredToken = fmt.Errorf("token is expired")
	ErrInvalidToken = fmt.Errorf("token is invalid")
)

type Payload struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(userID string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		SessionID: tokenID,
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (p *Payload) ValidExpiration() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
