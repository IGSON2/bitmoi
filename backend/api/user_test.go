package api

import (
	mockdb "bitmoi/backend/db/mock"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	mocktask "bitmoi/backend/worker/mock"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	if err := utilities.CheckPassword(m.password, actualParam.HashedPassword); err != nil {
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
	user, password := randomUser(t)
	hashed, err := utilities.HashPassword(password)

	storeCtrl := gomock.NewController(t)
	mockStore := mockdb.NewMockStore(storeCtrl)

	taskCtrl := gomock.NewController(t)
	mockTask := mocktask.NewMockTaskDistributor(taskCtrl)

	mockStore.EXPECT().GetAllParisInDB(gomock.Any()).Times(1)

	s := newTestServer(t, mockStore, mockTask)

	require.NoError(t, err)
	param := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			UserID:         user.UserID,
			Nickname:       user.Nickname,
			HashedPassword: hashed,
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
		Nickname: user.Nickname,
		Email:    user.Email,
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	httpReq, err := http.NewRequest("POST", "/user", bytes.NewReader(b))
	require.NoError(t, err)
	httpReq.Header.Set("Content-Type", "application/json")

	res, err := s.router.Test(httpReq)
	require.NoError(t, err)

	resBody, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Contains(t, string(resBody), "@")
}
