package api

import (
	db "bitmoi/db/sqlc"
	"bitmoi/token"
	"bitmoi/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

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
	testCount              = 5
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
}

var wg sync.WaitGroup

func TestSomePairs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	store := newTestStore(t)
	server := newTestServer(t, store, nil)

	ch := make(chan testResult, 2)
	wg.Add(testCount * 2)

	for i := 0; i < testCount; i++ {
		go testAnotherInterval(t, store, server, ch)
	}
	for i := 0; i < testCount*2; i++ {
		tr := <-ch
		t.Logf("result received name: %s", tr.intervalRes.Name)
		go testResponseWithRequest(t, tr.candleRes, tr.intervalRes, tr.intervalReq)
	}
	wg.Wait()
}

func testAnotherInterval(t *testing.T, store db.Store, s *Server, ch chan<- testResult) {
	defer func() {
		r := recover()
		if r != nil {
			t.Log("Sender goroutine terminated with error:", r)
			close(ch)
		}
	}()

	testCases := []struct {
		Name        string
		Path        string
		IntervalReq *AnotherIntervalRequest
		SetUpAuth   func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker)
	}{
		{
			Name: OKPractice4h,
			Path: "/practice",
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FourH,
				Mode:        practice,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
			},
		},
		{
			Name: OKCompetition4h,
			Path: "/competition",
			IntervalReq: &AnotherIntervalRequest{
				ReqInterval: db.FourH,
				Mode:        competition,
			},
			SetUpAuth: func(t *testing.T, request *http.Request, tokenMaker token.PasetoMaker) {
				addAuthrization(t, request, s.tokenMaker, authorizationTypeBearer, "igson", time.Minute)
			},
		},
	}

	for _, tc := range testCases {
		tn := new(TestName)
		t.Run(tc.Name, func(t *testing.T) {
			originOC := new(OnePairChart)
			intervalOC := new(OnePairChart)

			getCandleUrl := fmt.Sprintf("%s%s?names=%s", apiAddress, tc.Path, tn.Names)
			candleReq, err := http.NewRequest("GET", getCandleUrl, nil)
			candleReq.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			tc.SetUpAuth(t, candleReq, *s.tokenMaker)

			candleRes, err := http.DefaultClient.Do(candleReq)
			require.NoError(t, err)

			body, err := ioutil.ReadAll(candleRes.Body)

			require.NotNil(t, body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, originOC), string(body))
			require.NotNil(t, originOC.Identifier)

			tc.IntervalReq.Identifier = originOC.Identifier

			intervalReq, err := http.NewRequest("GET", apiAddress+encodeParams(tc.IntervalReq), nil)
			intervalReq.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			tc.SetUpAuth(t, intervalReq, *s.tokenMaker)

			intervalRes, err := http.DefaultClient.Do(intervalReq)
			require.NoError(t, err)

			body, err = ioutil.ReadAll(intervalRes.Body)

			require.NotNil(t, body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(body, intervalOC), string(body))
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
			}
			t.Logf("result sended testcase: %s, name: %s", tc.Name, intervalOC.Name)
		})
	}

}

func testResponseWithRequest(t *testing.T, candleRes, res *OnePairChart, req *AnotherIntervalRequest) {
	testIdentifier(t, res.Identifier, req.Identifier, req.ReqInterval)

	require.NotNil(t, candleRes)
	require.NotNil(t, candleRes.OneChart.PData)
	require.NotNil(t, candleRes.OneChart.VData)
	require.NotEmpty(t, candleRes.Identifier)
	require.NotNil(t, res)
	require.NotNil(t, res.OneChart.PData)
	require.NotNil(t, res.OneChart.VData)
	require.NotEmpty(t, res.Identifier)

	resTime := res.OneChart.PData[0].Time
	candleTime := candleRes.OneChart.PData[0].Time

	if req.Mode == practice {
		require.Equal(t, res.Name, candleRes.Name)
	}
	require.LessOrEqual(t, candleTime-resTime, db.GetIntervalStep(req.ReqInterval))
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

func encodeParams(intervalReq *AnotherIntervalRequest) string {
	params := url.Values{}
	params.Set("reqinterval", intervalReq.ReqInterval)
	params.Set("identifier", intervalReq.Identifier)
	params.Set("mode", intervalReq.Mode)
	return fmt.Sprintf("/interval?%s", params.Encode())
}
