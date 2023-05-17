package futureclient

import (
	"bitmoi/backend/utilities"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExchangeInfo(t *testing.T) {
	f, err := NewFutureClient(*utilities.GetConfig())
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
	f, err := NewFutureClient(*utilities.GetConfig())
	require.NoError(t, err)
	require.NotNil(t, f.Store)

	infos, err := f.getInfos(4, "h", "BTCUSDT")
	require.NoError(t, err)
	require.NotNil(t, infos)
}