package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPracRanks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)

	time1 := time.Now().Add(-3 * 24 * time.Hour)
	time2 := time.Now()
	rows, err := store.GetUserPracRankByPNL(context.Background(), GetUserPracRankByPNLParams{
		CreatedAt:   time1,
		CreatedAt_2: time2,
		Limit:       10,
		Offset:      0,
	})

	require.NoError(t, err)
	require.NotNil(t, rows)
	require.NotEmpty(t, rows)

	t.Logf("rows: %v", rows)
}
