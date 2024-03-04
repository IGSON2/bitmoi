package api

import (
	db "bitmoi/backend/db/sqlc"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRanks(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	s := newTestServer(t, newTestStore(t), nil)
	client := http.DefaultClient

	localAPI := fmt.Sprintf("http://localhost:%s", strings.Split(s.config.HTTPAddress, ":")[1])

	req, err := http.NewRequest("GET", localAPI+"/rank?mode=practice&category=pnl&start=24-03-01&end=24-03-04", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, res)

	rows := new(db.GetUserPracRankByPNLRow)
	b, err := io.ReadAll(res.Body)
	json.Unmarshal(b, rows)
	defer res.Body.Close()

	require.NotNil(t, rows)
	require.NoError(t, err)

	t.Log(rows)
}
