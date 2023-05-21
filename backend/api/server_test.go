package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"database/sql"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
)

func TestGetCandles(t *testing.T) {
	c := utilities.GetConfig("../../.")
	conn, err := sql.Open(c.DBDriver, c.DBSource)
	require.NoError(t, err)
	require.NotNil(t, conn)

	s := db.NewStore(conn)
	require.NotNil(t, s)

	server, err := NewServer(c, s)
	require.NoError(t, err)
	require.NotNil(t, server)

	req, err := http.NewRequest(http.MethodGet, "/test/?interval=4h&name=YGBUSDT", nil)
	require.NoError(t, err)

	res, err := server.router.Test(req)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestCalculateSeconds(t *testing.T) {
	s := 150 * 24 * time.Hour.Seconds()
	t.Log(s)
}
