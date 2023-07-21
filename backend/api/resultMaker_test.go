package api

import (
	"bitmoi/backend/utilities"
	"bitmoi/backend/utilities/common"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	localAPI = "http://localhost:5001"
)

func TestEntryPrice(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	s := newTestServer(t, newTestStore(t), nil)
	client := http.DefaultClient

	req, err := http.NewRequest("GET", localAPI+"/auth/competition", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	addAuthrization(t, req, s.tokenMaker, authorizationTypeBearer, masterID, 24*time.Hour)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, res)

	oc := new(OnePairChart)
	b, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	err = json.Unmarshal(b, oc)
	require.NoError(t, err)

	var info = new(utilities.IdentificationData)
	require.NoError(t, json.Unmarshal(utilities.DecryptByASE(oc.Identifier), info))

	var long = true

	compPostParam := ScoreRequest{
		Mode:        competition,
		UserId:      masterID,
		Name:        "STAGE 01",
		Stage:       1,
		IsLong:      &long,
		EntryPrice:  oc.EntryPrice,
		Quantity:    (defaultBalance * 50) / oc.EntryPrice,
		ProfitPrice: oc.EntryPrice + (oc.EntryPrice * 0.05),
		LossPrice:   oc.EntryPrice - (oc.EntryPrice * 0.0199),
		Leverage:    50,
		Balance:     defaultBalance,
		Identifier:  oc.Identifier,
		ScoreId:     "12345678901234",
		WaitingTerm: 1,
	}

	reqByte, err := json.Marshal(compPostParam)
	require.NoError(t, err)

	req2, err := http.NewRequest("POST", localAPI+"/auth/competition", bytes.NewReader(reqByte))
	require.NoError(t, err)
	req2.Header.Set("Content-Type", "application/json")
	addAuthrization(t, req2, s.tokenMaker, authorizationTypeBearer, masterID, 24*time.Hour)

	res2, err := client.Do(req2)
	require.NoError(t, err)
	require.NotNil(t, res2)

	resultBytes, err := ioutil.ReadAll(res2.Body)
	require.NoError(t, err)

	scoreResult := new(ScoreResponse)
	err = json.Unmarshal(resultBytes, scoreResult)
	require.NoError(t, err)

	decryptedOrigin := common.FloorDecimal(oc.OneChart.PData[0].Close / info.PriceFactor)
	require.Equal(t, scoreResult.OriginChart.PData[0].Close, decryptedOrigin)
}
