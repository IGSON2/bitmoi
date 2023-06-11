package api

import (
	"bitmoi/backend/utilities"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {

}

func randomUserRequest(t *testing.T) (*http.Request, string) {
	password := utilities.MakeRanString(6)
	hashed, err := utilities.HashPassword(password)
	require.NoError(t, err)

	b, err := json.Marshal(CreateUserRequest{
		UserID:   utilities.MakeRanString(5),
		Password: password,
		FullName: utilities.MakeRanString(10),
		Email:    utilities.MakeRanEmail(),
	})
	require.NoError(t, err)

	httpReq, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(b))
	require.NoError(t, err)
	httpReq.Header.Set("content-type", "application/json")

	return httpReq, hashed
}
