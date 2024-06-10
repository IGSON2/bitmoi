package api

import (
	"bitmoi/mail"
	"bitmoi/utilities"
	mocktask "bitmoi/worker/mock"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/mock/gomock"
	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/require"
)

func TestAsyncWorker(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config := utilities.GetConfig("../..")
	store := newTestStore(t)

	taskCtrl := gomock.NewController(t)
	mockTask := mocktask.NewMockTaskDistributor(taskCtrl)

	server := newTestServer(t, store, mockTask)

	processor := NewTestRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}, store, &mail.GmailSender{}, time.Minute)
	go processor.Start()
	for i := 0; i < 10; i++ {
		user := CreateUserRequest{
			UserID:   utilities.MakeRanString(6),
			Password: "secret123",
			Nickname: utilities.MakeRanString(4),
			Email:    utilities.MakeRanEmail(),
		}
		b, err := json.Marshal(user)
		require.NoError(t, err)
		req, err := http.NewRequest("POST", "/user", bytes.NewReader(b))
		require.NoError(t, err)
		req.Header.Set("content-type", "application/json")

		res, err := server.router.Test(req)
		require.NoError(t, err)
		body, err := ioutil.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode, string(body))
	}
	go server.Listen()
	time.Sleep(25 * time.Second)
}

func TestNotVerifiedLogin(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config := utilities.GetConfig("../..")
	store := newTestStore(t)

	taskCtrl := gomock.NewController(t)
	mockTask := mocktask.NewMockTaskDistributor(taskCtrl)

	server := newTestServer(t, store, mockTask)

	processor := NewTestRedisTaskProcessor(asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}, store, &mail.GmailSender{}, time.Minute)
	go processor.Start()

	login := LoginUserRequest{
		UserID:   "test12",
		Password: "secret",
	}
	b, err := json.Marshal(login)
	require.NoError(t, err)
	req, err := http.NewRequest("POST", "/user/login", bytes.NewReader(b))
	require.NoError(t, err)
	req.Header.Set("content-type", "application/json")

	res, err := server.router.Test(req)
	require.NoError(t, err)
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode, string(body))
	fmt.Println(string(body))

	go server.Listen()
	time.Sleep(25 * time.Second)
}
