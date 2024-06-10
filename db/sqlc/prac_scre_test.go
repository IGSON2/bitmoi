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

	time1, err := time.Parse("06-01-02", "24-05-06")
	require.NoError(t, err)

	time2, err := time.Parse("06-01-02", "24-05-07")
	require.NoError(t, err)

	rows, err := store.GetUserPracRankByPNL(context.Background(), GetUserPracRankByPNLParams{
		CreatedAt:   time1.Add(-9 * time.Hour),
		CreatedAt_2: time2.Add(-9 * time.Hour),
		Limit:       10,
		Offset:      0,
	})

	require.NoError(t, err)
	require.NotNil(t, rows)
	require.NotEmpty(t, rows)

	t.Logf("rows: %v", rows)
}
