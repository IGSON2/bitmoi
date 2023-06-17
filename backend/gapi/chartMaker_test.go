package gapi

import (
	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/gapi/pb"
	"bitmoi/backend/utilities"
	"context"
	"testing"

	"github.com/golang/mock/gomock"
)

var (
	user = utilities.MakeRanString(6)
)

func TestMakeChart(t *testing.T) {

	testCases := []struct {
		Name      string
		req       *pb.GetCandlesRequest
		buidStubs func(store *mockdb.MockStore)
		SetUpAuth func(t *testing.T, token string) context.Context
		CheckResp func(resp *pb.GetCandlesResponse)
	}{
		{
			Name: "OK_Practice",
			req: &pb.GetCandlesRequest{
				Names:  "",
				Mode:   practice,
				UserId: user,
			},
			buidStubs: func(store *mockdb.MockStore) {
				arg := db.Get1hCandlesParams{}
				rsp := []db.Candles1h{}
				store.EXPECT().Get1hCandles(gomock.Any(), arg).Times(1).Return(rsp)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			mockStore := mockdb.NewMockStore(storeCtrl)

			server := newTestServer(t, mockStore)
		})
	}

}
