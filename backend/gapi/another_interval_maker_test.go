package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

const (
	OKPractice5m           = "OK_Practice_5m"
	OKPractice15m          = "OK_Practice_15m"
	FailPracticeValidation = "Fail_Validation_Practice"
	OKCompetition5m        = "OK_Competition_5m"
	OKCompetition15m       = "OK_Competition_15m"
	FailCompetitionNoAuth  = "Fail_Competition_NoAuth"
)

var ()

type TestName struct {
	Names string
}

func (t *TestName) append(pair string) {
	if strings.Contains(strings.ToLower(pair), "stage") {
		log.Error().Err(fmt.Errorf("cannot append competition pariname: pair=%s", pair))
		return
	}
	if n := strings.Count(t.Names, ","); n == 9 {
		t.Names += pair
		return
	} else if n > 9 {
		log.Error().Err(fmt.Errorf("cannot append anymore: n=%d", n))
		return
	}
	t.Names += pair + ","
}

type testResult struct {
	testName    string
	candleRes   *pb.CandlesResponse
	intervalRes *pb.CandlesResponse
	intervalReq *pb.AnotherIntervalRequest
	err         error
}

func TestSomePairs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	oc := sync.Once{}
	oc.Do(
		func() {
			var err error
			store = newTestStore(t)
			server = newTestServer(t, store)
			client = newGRPCClient(t)
			pairs, err = store.GetAllParisInDB(context.Background())
			require.NoError(t, err)
		},
	)

	ch := make(chan testResult, len(pairs))

	for i := 0; i < len(pairs); i++ {
		go testAnotherInterval(t, store, server, client, ch)
	}
	for i := 0; i < len(pairs); i++ {
		tr := <-ch
		if tr.testName == FailCompetitionNoAuth {
			require.Error(t, tr.err)
		} else {
			go testResponseWithRequest(t, tr.candleRes, tr.intervalRes, tr.intervalReq, tr.err)
		}
	}
}

func testAnotherInterval(t *testing.T, store db.Store, s *Server, client pb.BitmoiClient, ch chan<- testResult) {
	tn := new(TestName)

	testCases := []struct {
		Name      string
		CandleReq *pb.CandlesRequest
		Req       *pb.AnotherIntervalRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
	}{
		{
			Name: OKPractice5m,
			CandleReq: &pb.CandlesRequest{
				Mode:   practice,
				UserId: "",
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        practice,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
		},
		{
			Name: OKPractice15m,
			CandleReq: &pb.CandlesRequest{
				Mode:   practice,
				UserId: "",
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        practice,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				return context.Background()
			},
		},
		{
			Name: OKCompetition5m,
			CandleReq: &pb.CandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        competition,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
		},
		{
			Name: OKCompetition15m,
			CandleReq: &pb.CandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        competition,
				UserId:      user,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
		},
		{
			Name: FailCompetitionNoAuth,
			CandleReq: &pb.CandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				Mode:   competition,
				UserId: "unauthorized user",
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
		},
		{
			Name: FailPracticeValidation,
			CandleReq: &pb.CandlesRequest{
				Mode:   competition,
				UserId: user,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: "Unsupported",
				Mode:        competition,
				UserId:      user,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := tc.SetUpAuth(t, s.tokenMaker)

			tc.CandleReq.Names = tn.Names
			candleRes, err := client.RequestCandles(ctx, tc.CandleReq)
			require.NoError(t, err)

			tn.append(candleRes.Name)
			if len(utilities.SplitPairnames(tn.Names)) > 10 {
				t.Log("The number of pairs has reached 10. End the test ")
				return
			}

			tc.Req.Identifier = candleRes.Identifier

			res, err := client.AnotherInterval(ctx, tc.Req)

			ch <- testResult{
				testName:    tc.Name,
				candleRes:   candleRes,
				intervalRes: res,
				intervalReq: tc.Req,
				err:         err,
			}
		})
	}

}

func testIdentifier(t *testing.T, resIden, reqIden, requestedInterval string) {
	resInfo := new(utilities.IdentificationData)
	json.Unmarshal(utilities.DecryptByASE(resIden), resInfo)

	reqInfo := new(utilities.IdentificationData)
	json.Unmarshal(utilities.DecryptByASE(reqIden), reqInfo)

	require.Equal(t, reqInfo.Interval, db.OneH)
	require.Equal(t, resInfo.Interval, requestedInterval)

	require.Equal(t, resInfo.Name, reqInfo.Name)
	require.Equal(t, resInfo.RefTimestamp, reqInfo.RefTimestamp)
	require.Equal(t, resInfo.PriceFactor, reqInfo.PriceFactor)
	require.Equal(t, resInfo.VolumeFactor, reqInfo.VolumeFactor)
	require.Equal(t, resInfo.TimeFactor, reqInfo.TimeFactor)
}

func testResponseWithRequest(t *testing.T, candleRes, res *pb.CandlesResponse, req *pb.AnotherIntervalRequest, err error) {
	require.NoError(t, err)
	testIdentifier(t, res.Identifier, req.Identifier, req.ReqInterval)

	require.NotNil(t, candleRes)
	require.NotNil(t, candleRes.OneChart.PData)
	require.NotNil(t, candleRes.OneChart.VData)
	require.NotEmpty(t, candleRes.Identifier)
	require.NotNil(t, res)
	require.NotNil(t, res.OneChart.PData)
	require.NotNil(t, res.OneChart.VData)
	require.NotEmpty(t, res.Identifier)

	resClose := res.OneChart.PData[0].Close
	candleClose := candleRes.OneChart.PData[0].Close

	resTime := res.OneChart.PData[0].Time
	candleTime := candleRes.OneChart.PData[0].Time

	require.Equal(t, res.Name, candleRes.Name)
	require.Equal(t, resTime, candleTime)
	require.GreaterOrEqual(t, float64(0.03), math.Abs(resClose-candleClose)/resClose) // 각 단위의 캔들의 종가의 차이가 3% 이하여야 함
	require.Equal(t, res.EntryTime, candleRes.EntryTime)
	require.GreaterOrEqual(t, float64(0.01), math.Abs(res.EntryPrice-candleRes.EntryPrice)/res.EntryPrice)
}
