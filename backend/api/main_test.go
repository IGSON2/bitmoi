package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	c := utilities.GetConfig("../../../.")

	s, err := NewServer(c, store)

	require.NoError(t, err)
	return s
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
