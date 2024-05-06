package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenRecomReq(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	s := newTestServer(t, newTestStore(t), nil)
	userID := fmt.Sprintf("%s@email.com", utilities.MakeRanString(10))
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

	reqData := CreateRecommendHistoryRequest{Code: "8502B59E96"}
	resData := createAuthorizedPostRequest(t, s, userID, "/auth/user/recommender", reqData)
	t.Log(string(resData))
}
