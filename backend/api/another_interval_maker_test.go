package api

import (
	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	apiAddress = "http://bitmoi.co.kr:5000/interval"
)

type getCandleParamsMatcher struct {
	arg      db.GetCandlesInterface
	interval string
}

func (m getCandleParamsMatcher) String() string {
	return fmt.Sprintf("matches parameter for request candles in each intervals (%s)", m.interval)
}

func (m getCandleParamsMatcher) Matches(x interface{}) bool {
	switch m.interval {
	case db.FiveM:
		actualIn, ok := x.(db.Get5mCandlesParams)
		if !ok {
			return false
		}

		actualExpected, ok := (m.arg).(*db.Get5mCandlesParams)
		if !ok {
			return false
		}

		return actualIn.Name == actualExpected.Name && actualIn.Limit == actualExpected.Limit &&
			actualIn.Time == actualExpected.Time && m.interval == actualIn.GetInterval()
	case db.FifM:
		actualIn, ok := x.(db.Get15mCandlesParams)
		if !ok {
			return false
		}

		actualExpected, ok := (m.arg).(*db.Get15mCandlesParams)
		if !ok {
			return false
		}

		return actualIn.Name == actualExpected.Name && actualIn.Limit == actualExpected.Limit &&
			actualIn.Time == actualExpected.Time && m.interval == actualIn.GetInterval()
	case db.OneH:
		actualIn, ok := x.(db.Get1hCandlesParams)
		if !ok {
			return false
		}

		actualExpected, ok := (m.arg).(*db.Get1hCandlesParams)
		if !ok {
			return false
		}

		return actualIn.Name == actualExpected.Name && actualIn.Limit == actualExpected.Limit &&
			actualIn.Time == actualExpected.Time && m.interval == actualIn.GetInterval()
	case db.FourH:
		actualIn, ok := x.(db.Get4hCandlesParams)
		if !ok {
			return false
		}

		actualExpected, ok := (m.arg).(*db.Get4hCandlesParams)
		if !ok {
			return false
		}

		return actualIn.Name == actualExpected.Name && actualIn.Limit == actualExpected.Limit &&
			actualIn.Time == actualExpected.Time && m.interval == actualIn.GetInterval()
	case db.OneD:
		actualIn, ok := x.(db.Get1dCandlesParams)
		if !ok {
			return false
		}

		actualExpected, ok := (m.arg).(*db.Get1dCandlesParams)
		if !ok {
			return false
		}

		return actualIn.Name == actualExpected.Name && actualIn.Limit == actualExpected.Limit &&
			actualIn.Time == actualExpected.Time && m.interval == actualIn.GetInterval()
	default:
		return false
	}
}

func newGetCandlesParamMatcher(arg db.GetCandlesInterface, interval string) getCandleParamsMatcher {
	return getCandleParamsMatcher{arg, interval}
}

type candleMatcher struct {
	Name        string
	RefTime     int64
	ReqInterval string
}

func (m candleMatcher) String() string {
	return fmt.Sprintf("matches result of candle request in each intervals (%s)", m.ReqInterval)
}

func (m candleMatcher) Matches(x interface{}) bool {
	fmt.Print(x)
	switch m.ReqInterval {
	case db.FiveM:
		actualIn, ok := x.(*Candles5mSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.ReqInterval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.ReqInterval)) && m.ReqInterval == actualIn.Interval()
	case db.FifM:
		actualIn, ok := x.(*Candles15mSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.ReqInterval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.ReqInterval)) && m.ReqInterval == actualIn.Interval()
	case db.OneH:
		actualIn, ok := x.(*Candles1hSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.ReqInterval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.ReqInterval)) && m.ReqInterval == actualIn.Interval()
	case db.FourH:
		actualIn, ok := x.(*Candles4hSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.ReqInterval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.ReqInterval)) && m.ReqInterval == actualIn.Interval()
	case db.OneD:
		actualIn, ok := x.(*Candles1dSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.ReqInterval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.ReqInterval)) && m.ReqInterval == actualIn.Interval()
	default:
		return false
	}
}

func newCandleMatcher(info utilities.IdentificationData) candleMatcher {
	return candleMatcher{
		Name:        info.Name,
		RefTime:     info.RefTimestamp,
		ReqInterval: info.Interval,
	}
}

func TestSomePairs(t *testing.T) {
	storeCtrl := gomock.NewController(t)
	defer storeCtrl.Finish()
	mockStore := mockdb.NewMockStore(storeCtrl)

	mockStore.EXPECT().GetAllParisInDB(gomock.Any()).Times(1).Return([]string{}, nil)

	server := newTestServer(t, mockStore, nil)

	token, _, err := server.tokenMaker.CreateToken("igson", time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	name := utilities.FindDiffPair(testGetAllPairs(t), []string{})
	min, max := testGetMinMaxTime(t, db.OneH, name)

	info := utilities.IdentificationData{
		Name:         name,
		Interval:     db.OneH,
		RefTimestamp: max - calculateRefTimestamp(max-min, name, db.OneH),
		PriceFactor:  1,
		VolumeFactor: 1,
		TimeFactor:   int64(86400 * (utilities.MakeRanInt(10950, 19000))),
	}

	identifier := utilities.EncrtpByASE(&info)

	// Test get another interval
	arg := db.Get4hCandlesParams{
		Name:  info.Name,
		Time:  info.RefTimestamp,
		Limit: oneTimeStageLoad,
	}
	mockStore.EXPECT().Get4hCandles(gomock.Any(), newGetCandlesParamMatcher(&arg, db.FourH)).Times(1).
		// Return(newCandleMatcher(info), nil)
		Return([]db.Candles4h{}, nil)

	reqURL := fmt.Sprintf("%s?%s", apiAddress, encodeParams(db.FourH, identifier, competition, 1))
	req, err := http.NewRequest("GET", reqURL, nil)
	require.NoError(t, err)
	req.Header.Set(authorizationHeaderKey, token)

	res, err := server.router.Test(req, int(time.Minute.Milliseconds()))
	require.NoError(t, err)
	require.NotNil(t, res)
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, res.Status, http.StatusOK, string(b))
}

func encodeParams(interval, identifier, mode string, stage int) string {
	params := url.Values{}
	params.Set("reqinterval", interval)
	params.Set("identifier", identifier)
	params.Set("mode", mode)
	params.Set("stage", strconv.Itoa(stage))
	return params.Encode()
}
