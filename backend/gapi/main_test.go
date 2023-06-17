package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	c := utilities.GetConfig("../../.")

	s, err := NewServer(c, store)
	require.NoError(t, err)
	return s
}

func newTestStore(t *testing.T) db.Store {
	c := utilities.GetConfig("../../.")
	conn, err := sql.Open(c.DBDriver, c.DBSource)
	require.NoError(t, err)
	return db.NewStore(conn)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func generateTestAccessToken(t *testing.T, tm *token.PasetoMaker) string {
	token, payload, err := tm.CreateToken(user, time.Minute)
	require.NoError(t, err)
	require.NotNil(t, payload)
	require.NotEqual(t, token, "")
	return token
}

func addAuthHeaderIntoContext(t *testing.T, token string) context.Context {
	require.NotEqual(t, token, "")
	MD := metadata.New(map[string]string{authorizationHeaderKey: token})
	require.NotNil(t, MD)
	return metadata.NewIncomingContext(context.Background(), MD)
}
