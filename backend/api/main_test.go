package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	apiAddress = "https://api.bitmoi.co.kr"
	masterID   = "igson"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	c := utilities.GetConfig("../../.")
	s, err := NewServer(c, store, taskDistributor)
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

func randomUser(t *testing.T) (user db.User, password string) {
	password = "secret123"
	hashed, err := utilities.HashPassword(password)
	require.NoError(t, err)
	var defaultTime time.Time
	user = db.User{
		UserID:            utilities.MakeRanString(8),
		HashedPassword:    sql.NullString{Valid: true, String: hashed},
		Nickname:          utilities.MakeRanString(10),
		Email:             utilities.MakeRanEmail(),
		PasswordChangedAt: defaultTime,
		CreatedAt:         defaultTime,
	}

	return user, password
}
