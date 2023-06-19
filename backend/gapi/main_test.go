package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	user = utilities.MakeRanString(6)
)

func newTestServer(t *testing.T, store db.Store) *Server {
	c := utilities.GetConfig("../../.")

	s, err := NewServer(c, store)
	require.NoError(t, err)

	go s.ListenGRPC()
	go s.ListenGRPCGateWay()
	return s
}

func newTestStore(t *testing.T) db.Store {
	c := utilities.GetConfig("../../.")
	conn, err := sql.Open(c.DBDriver, c.DBSource)
	require.NoError(t, err)
	return db.NewStore(conn)
}

func newGRPCClient(t *testing.T) pb.BitmoiClient {
	conn, err := grpc.Dial("localhost:6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewBitmoiClient(conn)
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

	bearerToken := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
	md := metadata.MD{
		authorizationHeaderKey: []string{
			bearerToken,
		},
	}
	// Client - Outgoing , Server - Incoming
	return metadata.NewOutgoingContext(context.Background(), md)
}
