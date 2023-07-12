package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	c := utilities.GetConfig("../../.")
	c.AccessTokenDuration = time.Minute
	s, err := NewServer(c, store, taskDistributor)
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

func testGetAllPairs(t *testing.T) []string {
	store := newTestStore(t)
	pairs, err := store.GetAllParisInDB(context.Background())
	require.NoError(t, err)
	return pairs
}

func testGetMinMaxTime(t *testing.T, interval, name string) (min, max int64) {
	store := newTestStore(t)
	min, max, err := store.SelectMinMaxTime(interval, name, context.Background())
	require.NoError(t, err)
	return min, max
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
