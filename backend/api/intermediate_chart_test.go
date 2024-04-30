package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCloseImdScore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	s := newTestServer(t, newTestStore(t), nil)
	userID := fmt.Sprintf("%s@email.com", utilities.MakeRanString(10))
	ScoreID := utilities.MakeRanInt(10000000, 99999999)

	ctx := context.Background()
	_, err := s.store.CreateUser(ctx, db.CreateUserParams{
		UserID:          userID,
		Nickname:        strings.Split(userID, "@")[0],
		Email:           userID,
		PracBalance:     1000000,
		CompBalance:     0,
		RecommenderCode: utilities.MakeRanString(8),
	})
	require.NoError(t, err)

	// defer func() {
	// _, err := s.store.DeleteUser(ctx, userID)
	// require.NoError(t, err)
	// _, err = s.store.DeletePracScore(ctx, db.DeletePracScoreParams{
	// UserID:  userID,
	// ScoreID: fmt.Sprintf("%d", ScoreID),
	// })
	// require.NoError(t, err)
	// }()

	token, _, err := s.tokenMaker.CreateToken(userID, time.Hour)
	require.NoError(t, err)

	url := "/basic/practice"
	httpReq, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)

	httpReq.Header.Add(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))

	res, err := s.router.Test(httpReq, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	resData, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	newOC := new(OnePairChart)
	err = json.Unmarshal(resData, newOC)
	require.NoError(t, err)

	require.NotNil(t, newOC.OneChart)
	require.NotEmpty(t, newOC.Identifier)
	require.NotEmpty(t, newOC.EntryPrice)

	user, err := s.store.GetUser(ctx, userID)
	require.NoError(t, err)

	isLong := true
	leverage := 1

	url = "/intermediate/init"
	imdInitData := ImdScoreRequest{
		Mode:        practice,
		UserId:      userID,
		ScoreId:     fmt.Sprintf("%d", ScoreID),
		Name:        newOC.Name,
		IsLong:      &isLong,
		EntryPrice:  newOC.EntryPrice,
		Quantity:    user.PracBalance / newOC.EntryPrice * 0.99,
		ProfitPrice: newOC.EntryPrice * 1.05,
		LossPrice:   newOC.EntryPrice * 0.95,
		Leverage:    int8(leverage),
		Identifier:  newOC.Identifier,
	}
	reqData, err := json.Marshal(imdInitData)
	require.NoError(t, err)

	httpReq, err = http.NewRequest("POST", url, bytes.NewReader(reqData))
	require.NoError(t, err)
	httpReq.Header.Add(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Add("Content-Type", "application/json")

	res, err = s.router.Test(httpReq, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	maxTimeStamp := newOC.OneChart.PData[0].Time + int64(utilities.MakeRanInt(1, 100))*db.GetIntervalStep(db.OneH)

	url = "/intermediate/close"
	imdCloseData := ImdCloseRequest{
		ImdScoreRequest: imdInitData,
		ReqInterval:     db.OneH,
		MinTimestamp:    maxTimeStamp - db.GetIntervalStep(db.OneH),
		MaxTimestamp:    maxTimeStamp,
	}

	reqData, err = json.Marshal(imdCloseData)
	require.NoError(t, err)

	httpReq, err = http.NewRequest("POST", url, bytes.NewReader(reqData))
	require.NoError(t, err)
	httpReq.Header.Add(authorizationHeaderKey, fmt.Sprintf("Bearer %s", token))
	httpReq.Header.Add("Content-Type", "application/json")

	res, err = s.router.Test(httpReq, -1)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

}
