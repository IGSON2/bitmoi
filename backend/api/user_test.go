package api

import (
	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	mocktask "bitmoi/backend/worker/mock"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type createuserMatcher struct {
	param    db.CreateUserTxParams
	password string
	user     db.User
}

func (m *createuserMatcher) Matches(x interface{}) bool {
	actualParam, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	if err := utilities.CheckPassword(m.password, actualParam.HashedPassword.String); err != nil {
		return false
	}

	m.param.HashedPassword = actualParam.HashedPassword

	if !reflect.DeepEqual(m.param.CreateUserParams, actualParam.CreateUserParams) {
		return false
	}

	err := actualParam.AfterCreate(m.user)
	return err == nil
}

func (m *createuserMatcher) String() string {
	return fmt.Sprintf("matches create user param. param:%v", m.param)
}

func newCreateUserMatcher(p db.CreateUserTxParams, password string, user db.User) *createuserMatcher {
	return &createuserMatcher{
		param:    p,
		password: password,
		user:     user,
	}
}

func TestCreateUser(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	user, password := randomUser(t)
	hashed, err := utilities.HashPassword(password)

	storeCtrl := gomock.NewController(t)
	defer storeCtrl.Finish()
	mockStore := mockdb.NewMockStore(storeCtrl)

	taskCtrl := gomock.NewController(t)
	defer taskCtrl.Finish()
	mockTask := mocktask.NewMockTaskDistributor(taskCtrl)

	mockStore.EXPECT().GetAllPairsInDB1H(gomock.Any()).Times(1)

	s := newTestServer(t, mockStore, mockTask)

	require.NoError(t, err)
	param := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			UserID:         user.UserID,
			Nickname:       user.Nickname,
			HashedPassword: sql.NullString{Valid: true, String: hashed},
			Email:          user.Email,
		}}
	mockStore.EXPECT().CreateUserTx(gomock.Any(), newCreateUserMatcher(param, password, user)).Times(1).
		Return(db.CreateUserTxResult{User: user}, nil)

	taskPayload := &worker.PayloadSendVerifyEmail{
		UserID: user.UserID,
	}
	mockTask.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).Times(1).
		Return(nil)

	req := CreateUserRequest{
		UserID:   user.UserID,
		Password: password,
		Nickname: user.Nickname.String,
		Email:    user.Email,
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	httpReq, err := http.NewRequest("POST", "/user", bytes.NewReader(b))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := s.router.Test(httpReq)
	require.NoError(t, err)

	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(resBody), "@")
}

func TestMoreInfo(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	userID, scoreID := "igson", "1691324536647"

	s := newTestServer(t, newTestStore(t), nil)

	url := fmt.Sprintf("/moreinfo?userid=%s&scoreid=%s", userID, scoreID)

	httpReq, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	res, err := s.router.Test(httpReq, 10000)
	require.NoError(t, err)
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	fmt.Println(string(b))
}

func TestGetLastUserID(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	s := newTestServer(t, newTestStore(t), nil)

	id, err := s.store.GetLastUserID(context.Background())
	require.NoError(t, err)

	t.Log(id)
}
