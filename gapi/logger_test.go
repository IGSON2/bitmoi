package gapi

import (
	"bitmoi/gapi/pb"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	if err := os.MkdirAll(filepath.Join("data", "grpc", "logs"), 0700); err != nil {
		log.Panic().Msgf("Err cannot make datadir: %v", err)
	}
	c := newGRPCClient(t)

	res, err := c.RequestCandles(context.Background(), &pb.CandlesRequest{Mode: practice})
	require.NoError(t, err)
	require.NotNil(t, res)
}
