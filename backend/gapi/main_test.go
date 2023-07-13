package gapi

import (
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
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

const (
	gapiAddress = "bitmoi.co.kr:6000"
	masterID    = "igson"
)

var (
	tm     *token.PasetoMaker
	client pb.BitmoiClient
)

func newTestPasetoMaker(t *testing.T) *token.PasetoMaker {
	c := utilities.GetConfig("../../.")
	tm, err := token.NewPasetoTokenMaker(c.SymmetricKey)
	require.NoError(t, err)
	return tm
}

func newGRPCClient(t *testing.T) pb.BitmoiClient {
	conn, err := grpc.Dial(gapiAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return pb.NewBitmoiClient(conn)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func generateTestAccessToken(t *testing.T, tm *token.PasetoMaker) string {
	token, payload, err := tm.CreateToken(masterID, time.Minute)
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
