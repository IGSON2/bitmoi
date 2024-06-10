package db

import (
	"bitmoi/utilities"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) Store {
	c := utilities.GetConfig("../../../.")
	conn, err := sql.Open(c.DBDriver, c.DBSource)
	require.NoError(t, err)
	return NewStore(conn)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
