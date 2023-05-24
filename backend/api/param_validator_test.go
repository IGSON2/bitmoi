package api

import (
	"bitmoi/backend/utilities"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	testCases := []struct {
		name     string
		params   validInsertRankParams
		expected func(t *testing.T, es *utilities.ErrorResponse)
	}{
		{
			name: "OK",
			params: validInsertRankParams{
				UserId:       "user",
				ScoreId:      "12345",
				DisplayName:  "DisplayedName",
				PhotoUrl:     "photo/url",
				Comment:      "Some Comment : 안녕하세요! 123",
				FinalBalance: 1928.9293,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse) {
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Equal(t, "", es.FailedField, m)
				require.Equal(t, "", es.Tag, m)
				require.Equal(t, "", es.Value, m)
			},
		},
		{
			name: "Fail_Missing_UserID",
			params: validInsertRankParams{
				UserId:       "",
				ScoreId:      "12345",
				DisplayName:  "DisplayedName",
				PhotoUrl:     "photo/url",
				Comment:      "안녕하세요~",
				FinalBalance: 1928.9293,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse) {
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Equal(t, "validInsertRankParams.UserID", es.FailedField, m)
				require.Equal(t, "required", es.Tag, m)
				require.Equal(t, "", es.Value, m)
			},
		},
		{
			name: "Fail_Missing_Comment",
			params: validInsertRankParams{
				UserId:       "user",
				ScoreId:      "12345",
				DisplayName:  "DisplayedName",
				PhotoUrl:     "photourl",
				Comment:      "",
				FinalBalance: 1002341.123,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse) {
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Equal(t, "validInsertRankParams.Comment", es.FailedField, m)
				require.Equal(t, "required", es.Tag, m)
				require.Equal(t, "", es.Value, m)
			},
		},
		{
			name: "Fail_Invalid_ScoreID",
			params: validInsertRankParams{
				UserId:       "user",
				ScoreId:      "스코어",
				DisplayName:  "DisplayedName",
				PhotoUrl:     "photo.url",
				Comment:      "Some Comment : 안녕하세요!",
				FinalBalance: 1928.9293,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse) {
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Equal(t, "validInsertRankParams.ScoreID", es.FailedField, m)
				require.Equal(t, "numeric", es.Tag, m)
				require.Equal(t, "스코어", es.Value, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errResponse := utilities.ValidateStruct(tc.params)
			if errResponse != nil {
				for _, es := range *errResponse {
					tc.expected(t, es)
				}
			}
		})
	}

}
