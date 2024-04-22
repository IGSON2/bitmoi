package db

import (
	"context"
	"testing"

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
