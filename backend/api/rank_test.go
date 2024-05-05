package api

import (
	db "bitmoi/backend/db/sqlc"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRanks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	s := newTestServer(t, newTestStore(t), nil)

	req, err := http.NewRequest("GET", "/basic/rank?mode=practice&category=pnl&start=24-04-29&end=24-05-05", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.router.Test(req, -1)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, err)
	require.NotNil(t, res)

	var rows []db.GetUserPracRankByPNLRow
	b, err := io.ReadAll(res.Body)
	json.Unmarshal(b, &rows)
	defer res.Body.Close()

	require.NotEmpty(t, rows)
	require.NoError(t, err)

	t.Log(rows)
}
