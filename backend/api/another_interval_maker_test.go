package api

import (
	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	identifier = "gy+Itl5buecRyHZXGCMTDensSXQEvlUUmbL/64M8Yf5+KRd4KUp/jDfFV6jT8ZaoofQ9XNQ4Bh0EcmEDVXa71O8lTaGaAFMFva3t8Y4JD34Dvj4ZzWpAjk1rPFswaTOs7Cw7BcSOVM7EnbeYWg=="
	requestURL = "http://bitmoi.co.kr:5000/interval"
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
	Name     string
	RefTime  int64
	Interval string
}

func (m candleMatcher) String() string {
	return fmt.Sprintf("matches result of candle request in each intervals (%s)", m.Interval)
}

func (m candleMatcher) Matches(x interface{}) bool {
	switch m.Interval {
	case db.FiveM:
		actualIn, ok := x.(*Candles5mSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.Interval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.Interval)) && m.Interval == actualIn.Interval()
	case db.FifM:
		actualIn, ok := x.(*Candles15mSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.Interval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.Interval)) && m.Interval == actualIn.Interval()
	case db.OneH:
		actualIn, ok := x.(*Candles1hSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.Interval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.Interval)) && m.Interval == actualIn.Interval()
	case db.FourH:
		actualIn, ok := x.(*Candles4hSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.Interval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.Interval)) && m.Interval == actualIn.Interval()
	case db.OneD:
		actualIn, ok := x.(*Candles1dSlice)
		if !ok {
			return false
		}

		actualRefTime := actualIn.InitCandleData().PData[0].Time

		return actualIn.Name() == m.Name && actualIn.Interval() == m.Interval &&
			(m.RefTime-actualRefTime < db.CalculateSeconds(m.Interval)) && m.Interval == actualIn.Interval()
	default:
		return false
	}
}

func newCandleMatcher(info utilities.IdentificationData) candleMatcher {
	return candleMatcher{
		Name:     info.Name,
		RefTime:  info.RefTimestamp,
		Interval: info.Interval,
	}
}

func TestSomePairs(t *testing.T) {
	storeCtrl := gomock.NewController(t)
	defer storeCtrl.Finish()
	mockStore := mockdb.NewMockStore(storeCtrl)

	server := newTestServer(t, mockStore, nil)

	info := new(utilities.IdentificationData)
	err := json.Unmarshal(utilities.DecryptByASE(identifier), info)
	require.NoError(t, err)

	anotherIntervalParam := AnotherIntervalRequest{
		ReqInterval: db.FourH,
		Identifier:  identifier,
		Mode:        practice,
		Stage:       1,
	}

	// Test get another interval
	arg := db.Get4hCandlesParams{
		Name:  info.Name,
		Time:  info.RefTimestamp,
		Limit: oneTimeStageLoad,
	}
	mockStore.EXPECT().Get4hCandles(gomock.Any(), newGetCandlesParamMatcher(&arg, db.FourH)).
		Times(1).Return(newCandleMatcher(*info), nil)

	req := http.NewRequest("GET", "/interval", body io.Reader)
	server.router.Test()

}
