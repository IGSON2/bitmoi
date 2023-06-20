package gapi

import (
	"bitmoi/backend/gapi/pb"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	if err := os.MkdirAll(filepath.Join("data", "grpc", "logs"), 0700); err != nil {
		log.Panic().Msgf("Err cannot make datadir: %v", err)
	}
	c := newGRPCClient(t)

	res, err := c.RequestCandles(context.Background(), &pb.GetCandlesRequest{Mode: practice})
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestLogFormat(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Any("message", "HI").Msg("Test message")
}
