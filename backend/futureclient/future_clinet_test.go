package futureclient

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExchangeInfo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	f, err := NewFutureClient(utilities.GetConfig("../../."))
	require.NoError(t, err)
	require.NotNil(t, f.Store)
	info, err := f.Client.NewExchangeInfoService().Do(context.Background())
	require.NoError(t, err)
	require.NotNil(t, info)
	var symbols []string
	for _, s := range info.Symbols {
		symbols = append(symbols, s.Symbol)
	}
	require.Greater(t, len(symbols), 0)
}

func TestGetInfo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	f, err := NewFutureClient(utilities.GetConfig("../../."))
	require.NoError(t, err)
	require.NotNil(t, f.Store)

	var cnt int
	err = f.StoreCandles(db.FiveM, "XRPUSDT", 1686355200000, true, &cnt)
	require.NoError(t, err)
}
