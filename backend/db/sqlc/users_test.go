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
	zone, err := time.LoadLocation("Asia/Seoul")
	require.NoError(t, err)

	// zone = time.UTC
	now := time.Now().In(zone)

	_, err = store.UpdateUserLastAccessedAt(context.Background(), UpdateUserLastAccessedAtParams{
		LastAccessedAt: sql.NullTime{Time: now.Add(9 * time.Hour), Valid: true},
		UserID:         userID,
	})
	require.NoError(t, err)

	userA, err := store.GetUser(context.Background(), userID)
	require.NoError(t, err)
	require.WithinDuration(t, now, userA.LastAccessedAt.Time, time.Second*5)

	compareTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour()-1, 0, 0, 0, time.Local)
	require.True(t, compareTime.Before(userA.LastAccessedAt.Time))
	t.Log(compareTime, userA.LastAccessedAt.Time)
}
