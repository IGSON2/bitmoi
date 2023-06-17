package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	symmetrickey = "00000000000000000000000000000000"
	userID       = "user123"
)

func TestPasetoMaker(t *testing.T) {
	p, err := NewPasetoTokenMaker(symmetrickey)
	require.NoError(t, err)

	duration := time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := p.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEqual(t, token, "")
	require.NotNil(t, payload)
	require.NotZero(t, payload.SessionID.ID())
	require.NotEqual(t, payload.UserID, "")
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiredAt, time.Second)

	payload2, err := p.VerifyToken(token)
	require.NoError(t, err)
	require.NotNil(t, payload2)
	require.Equal(t, payload.SessionID, payload2.SessionID)
	require.Equal(t, payload.UserID, payload2.UserID)
	require.WithinDuration(t, payload2.IssuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, payload2.ExpiredAt, payload.ExpiredAt, time.Second)
}

func TestTokenExpire(t *testing.T) {
	p, err := NewPasetoTokenMaker(symmetrickey)
	require.NoError(t, err)

	duration := -1 * time.Minute

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(duration)

	token, payload, err := p.CreateToken(userID, duration)
	require.NoError(t, err)
	require.NotEqual(t, token, "")
	require.NotNil(t, payload)
	require.NotZero(t, payload.SessionID.ID())
	require.NotEqual(t, payload.UserID, "")
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiresAt, payload.ExpiredAt, time.Second)

	payload2, err := p.VerifyToken(token)
	require.Nil(t, payload2)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
}
