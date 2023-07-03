package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store) *Server {
	c := utilities.GetConfig("../../.")
	c.AccessTokenDuration = time.Minute
	s, err := NewServer(c, store)
	s.taskDistributor = NewRedisTaskDistributor(asynq.RedisClientOpt{Addr: c.RedisAddress})
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
