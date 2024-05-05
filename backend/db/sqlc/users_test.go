package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGetUsers(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)

	rows, err := store.GetUsers(context.Background(), GetUsersParams{
		Limit:  1000,
		Offset: 0,
	})

	require.NoError(t, err)
	t.Logf("rows: %v", rows)
}

func TestLastaccessedAt(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)

	userID := "coactadmin"

	now := time.Now()

	_, err := store.UpdateUserLastAccessedAt(context.Background(), UpdateUserLastAccessedAtParams{
		LastAccessedAt: sql.NullTime{Time: now, Valid: true},
		UserID:         userID,
	})
	require.NoError(t, err)

	userA, err := store.GetUser(context.Background(), userID)
	require.NoError(t, err)

	crZ, _ := userA.CreatedAt.Zone()
	laZ, _ := userA.LastAccessedAt.Time.Zone()
	nowZ, _ := now.Zone()
	t.Logf("now: %v, created_at: %v, last_accessed_at: %v", nowZ, crZ, laZ)
	t.Logf("now: %v, last_accessed_at: %v", now, userA.LastAccessedAt.Time)
}
