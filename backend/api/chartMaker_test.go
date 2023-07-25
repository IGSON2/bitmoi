package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestMakeChart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)
	s := newTestServer(t, store, nil)

	testCases := []struct {
		Name      string
		Path      string
		SetUpAuth func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
		CheckResp func(resp *http.Response)
	}{
		{
			Name:      "GetPracticeChart",
			Path:      "/practice",
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {},
			CheckResp: func(resp *http.Response) {
				oc := new(OnePairChart)

				require.Equal(t, fiber.StatusOK, resp.StatusCode)
				b, err := ioutil.ReadAll(resp.Body)

				require.NotNil(t, b)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(b, oc))
				require.NotEmpty(t, oc.Identifier)
				require.NotContains(t, oc.Identifier, " ")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(oc.Identifier), info))

				require.NotEmpty(t, oc.EntryTime)
				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, oc.BtcRatio, float64(0))

				require.NotNil(t, oc.OneChart.PData)
				require.NotNil(t, oc.OneChart.VData)

				require.Contains(t, info.Name, "USDT")
				require.Equal(t, info.Interval, db.OneH)
				require.Equal(t, info.PriceFactor, float64(0))
				require.Equal(t, info.VolumeFactor, float64(0))
				require.Equal(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name: "GetCompetitonChart",
			Path: "/competition",
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, s.tokenMaker, authorizationTypeBearer, masterID, time.Minute)
			},
			CheckResp: func(resp *http.Response) {
				oc := new(OnePairChart)

				require.Equal(t, fiber.StatusOK, resp.StatusCode)
				b, err := ioutil.ReadAll(resp.Body)

				require.NotNil(t, b)
				require.NoError(t, err)
				require.NoError(t, json.Unmarshal(b, oc))
				require.NotEmpty(t, oc.Identifier)
				require.NotContains(t, oc.Identifier, " ")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(oc.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, oc.BtcRatio, float64(0))
				require.NotEmpty(t, oc.EntryTime)

				require.NotNil(t, oc.OneChart.PData)
				require.NotNil(t, oc.OneChart.VData)

				require.Contains(t, info.Name, "USDT")
				require.Equal(t, info.Interval, db.OneH)
				require.Greater(t, info.PriceFactor, float64(0))
				require.Greater(t, info.VolumeFactor, float64(0), fmt.Sprintf("name: %s, refTime: %d", info.Name, info.RefTimestamp))
				require.Greater(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name: "CompetitionUnAuthorized",
			Path: "/competition",
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
			CheckResp: func(resp *http.Response) {
				require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			client := http.DefaultClient

			req, err := http.NewRequest("GET", apiAddress+tc.Path, nil)
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			tc.SetUpAuth(t, req, *s.tokenMaker)

			res, err := client.Do(req)
			require.NoError(t, err)

			tc.CheckResp(res)
		})
	}

}

func BenchmarkHTTPMakeChart(b *testing.B) {
	t := new(testing.T)
	req := makeTestRequest(t)
	c := &http.Client{}
	for i := 0; i < b.N; i++ {
		c.Do(req)
	}
}

func BenchmarkHTTPPostScore(b *testing.B) {
	t := new(testing.T)
	req := makeTestScoreRequest(t)
	c := &http.Client{}
	for i := 0; i < b.N; i++ {
		c.Do(req)
	}
}

func makeTestRequest(t *testing.T) *http.Request {
	b, err := json.Marshal(CandlesRequest{Names: ""})
	require.NoError(t, err)
	req, err := http.NewRequest("GET", "http://localhost:4000/practice", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func makeTestScoreRequest(t *testing.T) *http.Request {
	long := true
	b, err := json.Marshal(ScoreRequest{
		Mode:        practice,
		UserId:      "user",
		Name:        "BTCUSDT",
		Stage:       1,
		IsLong:      &long,
		EntryPrice:  1639.31,
		Quantity:    10,
		ProfitPrice: 1700,
		LossPrice:   1600,
		Leverage:    17,
		Balance:     1000,
		Identifier:  "ALYJ/z8Bnb4k2TwsZlSr1KAcxn/Km0IYrTKE3fayRnKvKCE2V4BiXe+el4m6g0j2QnBG13nziUjQ52v00pf4CoruyccVqkubqM0IEBL9jXMdz6VwtibVkVhxIlJMNwwQH3juPDGziIYw48Jq7g==",
		ScoreId:     "abc",
		WaitingTerm: 1,
	})
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "http://localhost:7001/v1/score", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}
