package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoTokenMaker(symmetricKey string) (*PasetoMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size : must be exactly %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (p *PasetoMaker) CreateToken(userID string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := p.paseto.Encrypt(p.symmetricKey, payload, nil)
	return token, payload, err
}

func (p *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := new(Payload)

	if err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil); err != nil {
		return nil, fmt.Errorf("%w %w", ErrInvalidToken, err)
	}

	if err := payload.ValidExpiration(); err != nil {
		return nil, err
	}

	return payload, nil
}
