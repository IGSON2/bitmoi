package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeChart(t *testing.T) {
	tm := newTestPasetoMaker(t)
	client := newGRPCClient(t)

	testCases := []struct {
		Name      string
		req       *pb.CandlesRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
		CheckResp func(res *pb.CandlesResponse, err error)
	}{
		{
			Name: "OK_Practice",
			req: &pb.CandlesRequest{
				Names:  "",
				Mode:   practice,
				UserId: masterID,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.CandlesResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(res.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, res.BtcRatio, float64(0))
				require.Greater(t, res.EntryTime, "")

				require.NotNil(t, res.OneChart.PData)
				require.NotNil(t, res.OneChart.VData)

				require.Contains(t, res.Name, "USDT")
				require.Equal(t, info.Interval, db.OneH)
				require.Equal(t, info.PriceFactor, float64(0))
				require.Equal(t, info.VolumeFactor, float64(0))
				require.Equal(t, info.TimeFactor, int64(0))
			},
		},
		{
			Name: "OK_Competition",
			req: &pb.CandlesRequest{
				Names:  "",
				Mode:   competition,
				UserId: masterID,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
			CheckResp: func(res *pb.CandlesResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, res.Identifier, "")

				info := new(utilities.IdentificationData)
				require.NoError(t, json.Unmarshal(utilities.DecryptByASE(res.Identifier), info))

				require.Greater(t, info.RefTimestamp, int64(0))
				require.Greater(t, res.BtcRatio, float64(0))
				require.Greater(t, res.EntryTime, "")

				require.NotNil(t, res.OneChart.PData)
				require.NotNil(t, res.OneChart.VData)

				require.Contains(t, res.Name, "STAGE")
				require.Equal(t, info.Interval, db.OneH)
				require.Greater(t, info.PriceFactor, float64(0), fmt.Sprintf("name:%s, ref:%d", info.Name, info.RefTimestamp))
				require.Greater(t, info.VolumeFactor, float64(0), fmt.Sprintf("name:%s, ref:%d", info.Name, info.RefTimestamp))
				require.Greater(t, info.TimeFactor, int64(0), fmt.Sprintf("name:%s, ref:%d", info.Name, info.RefTimestamp))
			},
		},
		{
			Name: "No_Auth_Competition",
			req: &pb.CandlesRequest{
				Names:  "",
				Mode:   competition,
				UserId: "unauthorized user",
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.CandlesResponse, err error) {
				require.Error(t, err)
			},
		},
		{
			Name: "Fail_Validation_Practice",
			req: &pb.CandlesRequest{
				Names:  "",
				Mode:   "Unsupported",
				UserId: masterID,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
			CheckResp: func(res *pb.CandlesResponse, err error) {
				t.Log(err)
				require.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := tc.SetUpAuth(t, tm)
			res, err := client.RequestCandles(ctx, tc.req)
			tc.CheckResp(res, err)
		})
	}

}

func BenchmarkMakeChart(b *testing.B) {
	t := new(testing.T)
	c := newGRPCClient(t)
	for i := 0; i < b.N; i++ {
		c.RequestCandles(context.Background(), &pb.CandlesRequest{Mode: practice})
	}
}

func BenchmarkPostScore(b *testing.B) {
	t := new(testing.T)
	c := newGRPCClient(t)
	for i := 0; i < b.N; i++ {
		c.PostScore(context.Background(), &pb.ScoreRequest{
			Mode:        practice,
			UserId:      "user",
			Name:        "BTCUSDT",
			Stage:       1,
			IsLong:      true,
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
	}
}

func BenchmarkGateWayMakeChart(b *testing.B) {
	t := new(testing.T)
	req := makeTestGateWayCandleRequest(t)
	c := &http.Client{}
	for i := 0; i < b.N; i++ {
		c.Do(req)
	}
}

func BenchmarkGateWayPostScore(b *testing.B) {
	t := new(testing.T)
	req := makeTestGateWayScoreRequest(t)
	c := &http.Client{}
	for i := 0; i < b.N; i++ {
		c.Do(req)
	}
}

func makeTestGateWayCandleRequest(t *testing.T) *http.Request {
	b, err := json.Marshal(pb.CandlesRequest{Mode: practice})
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "http://localhost:7001/v1/candles", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func makeTestGateWayScoreRequest(t *testing.T) *http.Request {
	b, err := json.Marshal(pb.ScoreRequest{
		Mode:        practice,
		UserId:      "user",
		Name:        "BTCUSDT",
		Stage:       1,
		IsLong:      true,
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
