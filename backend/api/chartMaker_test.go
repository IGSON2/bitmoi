package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestMakeChart(t *testing.T) {
	var candles CandlesInterface
	testCases := []struct {
		Name      string
		Path      string
		Method    string
		Body      interface{}
		BuildStub func(s *mockdb.MockStore)
		CheckResp func(resp *http.Response)
	}{
		{
			Name:   "GetPracticeChart",
			Path:   "/practice",
			Method: http.MethodGet,
			Body: ChartRequestQuery{
				Names:    "",
				Interval: db.FourH,
			},
			BuildStub: func(s *mockdb.MockStore) {
				s.EXPECT().Get4hMinMaxTime(gomock.Any(), gomock.Any()).Times(1).Return(db.Get4hMinMaxTimeRow{}, nil)
				s.EXPECT().Get4hCandles(gomock.Any(), gomock.Any()).Times(1).Return(candles, nil)
			},
			CheckResp: func(resp *http.Response) {
				require.Equal(t, http.StatusOK, resp.Status)
				require.NotNil(t, resp.Body)
			},
		},
	}

	ctr := gomock.NewController(t)
	store := mockdb.NewMockStore(ctr)
	server := newTestServer(t, store)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			b, err := json.Marshal(tc.Body)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.Method, tc.Path, bytes.NewReader(b))
			require.NoError(t, err)

			res, err := server.router.Test(req)
			require.NoError(t, err)

			tc.CheckResp(res)
		})
	}

}
