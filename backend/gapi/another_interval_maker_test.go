package gapi

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

const (
	OKPractice4h           = "OK_Practice_4h"
	OKPractice15m          = "OK_Practice_15m"
	FailPracticeValidation = "Fail_Validation_Practice"
	OKCompetition4h        = "OK_Competition_4h"
	OKCompetition15m       = "OK_Competition_15m"
	FailCompetitionNoAuth  = "Fail_Competition_NoAuth"
	testCount              = 15
)

type TestName struct {
	Names string
}

var wg sync.WaitGroup

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
	oc := sync.Once{}
	oc.Do(
		func() {
			var err error
			tm = newTestPasetoMaker(t)
			client = newGRPCClient(t)
			require.NoError(t, err)
		},
	)

	wg.Add(testCount)

	ch := make(chan testResult)

	for i := 0; i < testCount; i++ {
		go testAnotherInterval(t, tm, client, ch)
	}
	for i := 0; i < testCount; i++ {
		tr := <-ch
		go testResponseWithRequest(t, tr.candleRes, tr.intervalRes, tr.intervalReq, tr.err)
	}
	wg.Wait()
}

func testAnotherInterval(t *testing.T, tm *token.PasetoMaker, client pb.BitmoiClient, ch chan<- testResult) {
	tn := new(TestName)

	testCases := []struct {
		Name      string
		CandleReq *pb.CandlesRequest
		Req       *pb.AnotherIntervalRequest
		SetUpAuth func(t *testing.T, tm *token.PasetoMaker) context.Context
	}{
		// {
		// 	Name: OKPractice4h,
		// 	CandleReq: &pb.CandlesRequest{
		// 		Mode:   practice,
		// 		UserId: "",
		// 	},
		// 	Req: &pb.AnotherIntervalRequest{
		// 		ReqInterval: db.FourH,
		// 		Mode:        practice,
		// 		UserId:      masterID,
		// 		Stage:       1,
		// 	},
		// 	SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
		// 		return context.Background()
		// 	},
		// },
		{
			Name: OKCompetition4h,
			CandleReq: &pb.CandlesRequest{
				Mode:   competition,
				UserId: masterID,
			},
			Req: &pb.AnotherIntervalRequest{
				ReqInterval: db.FourH,
				Mode:        competition,
				UserId:      masterID,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, tm *token.PasetoMaker) context.Context {
				token := generateTestAccessToken(t, tm)
				return addAuthHeaderIntoContext(t, token)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := tc.SetUpAuth(t, tm)

			tc.CandleReq.Names = tn.Names
			candleRes, err := client.RequestCandles(ctx, tc.CandleReq)
			require.NoError(t, err)
			if candleRes.Identifier == "" {
				t.Error("identifier is nil")
			}

			tn.append(candleRes.Name)
			if len(utilities.SplitPairnames(tn.Names)) > 10 {
				t.Log("The number of pairs has reached 10. End the test ")
				return
			}

			tc.Req.Identifier = candleRes.Identifier

			res, err := client.AnotherInterval(ctx, tc.Req)
			fmt.Printf("name:%s, len:%d err:%v\n", res.Name, len(res.OneChart.PData), err)
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
	require.NotContains(t, candleRes.Identifier, " ")

	require.NotNil(t, res)
	require.NotNil(t, res.OneChart.PData)
	require.NotNil(t, res.OneChart.VData)
	require.NotEmpty(t, res.Identifier)

	resTime := res.OneChart.PData[0].Time
	candleTime := candleRes.OneChart.PData[0].Time

	if req.Mode == practice {
		require.Equal(t, res.Name, candleRes.Name)
	}
	require.LessOrEqual(t, candleTime-resTime, db.CalculateSeconds(req.ReqInterval))
	wg.Done()
}
