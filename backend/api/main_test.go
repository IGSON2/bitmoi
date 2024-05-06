package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func createAuthorizedGetRequest(t *testing.T, s *Server, userID, url string) []byte {

	token, _, err := s.tokenMaker.CreateToken(userID, time.Hour)
	require.NoError(t, err)

	httpReq, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	httpReq.Header.Add(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Add("Content-Type", "application/json")

	res, err := s.router.Test(httpReq, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return resData
}

func createAuthorizedPostRequest(t *testing.T, s *Server, userID, url string, data any) []byte {

	reqData, err := json.Marshal(data)
	require.NoError(t, err)

	token, _, err := s.tokenMaker.CreateToken(userID, time.Hour)
	require.NoError(t, err)

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(reqData))
	require.NoError(t, err)

	httpReq.Header.Add(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Add("Content-Type", "application/json")

	res, err := s.router.Test(httpReq, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return resData
}
