package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

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
	candleRes   *OnePairChart
	intervalRes *OnePairChart
	intervalReq *AnotherIntervalRequest
	err         error
}

var wg sync.WaitGroup

func TestSomePairs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)
	server := newTestServer(t, store)
	pairs, err := store.GetAllParisInDB(context.Background())
	require.NoError(t, err)

	ch := make(chan testResult)
	wg.Add(len(pairs))

	for i := 0; i < len(pairs); i++ {
		go testAnotherInterval(t, store, server, ch)
	}
	for i := 0; i < len(pairs); i++ {
		tr := <-ch
		go testResponseWithRequest(t, tr.candleRes, tr.intervalRes, tr.intervalReq, tr.err)
	}
	wg.Wait()
}

func testAnotherInterval(t *testing.T, store db.Store, s *Server, ch chan<- testResult) {
	tn := new(TestName)

	testCases := []struct {
		Name        string
		CandleReq   *CandlesRequest
		Path        string
		IntervalReq *AnotherIntervalRequest
		SetUpAuth   func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
	}{
		{
			Name:      OKPractice5m,
			Path:      "/practice",
			CandleReq: &CandlesRequest{},
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        practice,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
		},
		{
			Name:      OKPractice15m,
			Path:      "/practice",
			CandleReq: &CandlesRequest{},
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        practice,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
		},
		{
			Name:      OKCompetition5m,
			Path:      "/auth/competition",
			CandleReq: &CandlesRequest{},
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FiveM,
				Mode:        competition,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, s.tokenMaker, authorizationTypeBearer, "bknuw", time.Minute)
			},
		},
		{
			Name:      OKCompetition15m,
			Path:      "/auth/competition",
			CandleReq: &CandlesRequest{},
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FifM,
				Mode:        competition,
				Stage:       1,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, s.tokenMaker, authorizationTypeBearer, "bknuw", time.Minute)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			originOC := new(OnePairChart)
			intervalOC := new(OnePairChart)

			b, err := json.Marshal(tc.CandleReq)
			require.NoError(t, err)

			candleReq, err := http.NewRequest("GET", tc.Path, bytes.NewReader(b))
			candleReq.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			tc.SetUpAuth(t, candleReq, *s.tokenMaker)

			candleRes, err := s.router.Test(candleReq, 60000)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(candleRes.Body)

			require.NotNil(t, body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, originOC))
			require.NotNil(t, originOC.Identifier)

			tc.CandleReq.Names = tn.Names
			tc.IntervalReq.Identifier = originOC.Identifier

			b, err = json.Marshal(tc.IntervalReq)
			require.NoError(t, err)

			intervalReq, err := http.NewRequest("GET", "/interval", bytes.NewReader(b))
			intervalReq.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			tc.SetUpAuth(t, intervalReq, *s.tokenMaker)

			intervalRes, err := s.router.Test(intervalReq, 60000)
			require.NoError(t, err)

			body, err = ioutil.ReadAll(intervalRes.Body)

			require.NotNil(t, body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, intervalOC), fmt.Sprintf("%s", body))
			require.NotNil(t, intervalOC.Identifier)

			tn.append(intervalOC.Name)
			if len(utilities.SplitPairnames(tn.Names)) > 10 {
				t.Log("The number of pairs has reached 10. End the test ")
				return
			}

			ch <- testResult{
				testName:    tc.Name,
				candleRes:   originOC,
				intervalRes: intervalOC,
				intervalReq: tc.IntervalReq,
				err:         err,
			}
		})
	}

}

func testResponseWithRequest(t *testing.T, candleRes, res *OnePairChart, req *AnotherIntervalRequest, err error) {
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
	require.GreaterOrEqual(t, float64(0.02), math.Abs(resClose-candleClose)/resClose) // 각 단위의 캔들의 종가의 차이가 2% 이하여야 함
	require.Equal(t, res.EntryTime, candleRes.EntryTime)
	require.GreaterOrEqual(t, float64(0.02), math.Abs(res.EntryPrice-candleRes.EntryPrice)/res.EntryPrice)
	wg.Done()
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
