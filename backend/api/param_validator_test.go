package api

import (
	"bitmoi/backend/utilities"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderstructValidation(t *testing.T) {
	var (
		// long  = true
		short = false
	)
	testCases := []struct {
		name     string
		params   ScoreRequest
		expected func(t *testing.T, es *utilities.ErrorResponse, o ScoreRequest, i int)
	}{
		{
			name: "OK",
			params: ScoreRequest{
				Mode:        competition,
				UserId:      "user",
				Name:        "A",
				Stage:       1,
				IsLong:      &short,
				EntryPrice:  100.123,
				Quantity:    0.123,
				ProfitPrice: 110.2521,
				LossPrice:   90.2817,
				Leverage:    25,
				Balance:     1000,
				Identifier:  "asdf",
				ScoreId:     "12345",
				WaitingTerm: 1,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o ScoreRequest, i int) {
				require.Nil(t, es)
			},
		},
		{
			name: "Fail_Missing_Fields",
			params: ScoreRequest{
				Mode:        "",
				UserId:      "",
				Name:        "",
				Stage:       -1,
				IsLong:      nil,
				EntryPrice:  -1,
				Quantity:    -1,
				ProfitPrice: -1,
				LossPrice:   -1,
				Leverage:    0,
				Balance:     -1,
				Identifier:  "",
				ScoreId:     "",
				WaitingTerm: 0,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o ScoreRequest, i int) {
				r := reflect.TypeOf(o).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(o).Field(i).Interface(), es.Value, m)
			},
		},
		{
			name: "Fail_Boundary_Value",
			params: ScoreRequest{
				Mode:        "",
				UserId:      "",
				Name:        "",
				Stage:       0,
				IsLong:      nil,
				EntryPrice:  0,
				Quantity:    0,
				ProfitPrice: 0,
				LossPrice:   0,
				Leverage:    0,
				Balance:     0,
				Identifier:  "",
				ScoreId:     "",
				WaitingTerm: 0,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o ScoreRequest, i int) {
				r := reflect.TypeOf(o).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(o).Field(i).Interface(), es.Value, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errResponse := utilities.ValidateStruct(tc.params)
			if errResponse != nil {
				for i, es := range *errResponse {
					tc.expected(t, es, tc.params, i)
				}
			}
		})
	}

}

func TestRankInsertRequestValidation(t *testing.T) {
	testCases := []struct {
		name     string
		params   RankInsertRequest
		expected func(t *testing.T, es *utilities.ErrorResponse, req RankInsertRequest, i int)
	}{
		{
			name: "OK",
			params: RankInsertRequest{
				ScoreId: "123",
				Comment: "comment",
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req RankInsertRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.NotContains(t, es.FailedField, r.Name, m)
				require.NotContains(t, r.Tag, es.Tag, m)
				require.NotEqual(t, reflect.ValueOf(req).Field(i), es.Value, m)
			},
		},
		{
			name: "Fail_Missing_Fields",
			params: RankInsertRequest{
				ScoreId: "score",
				Comment: "",
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req RankInsertRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				if r.Name == "Comment" || r.Name == "Nickname" {
					return
				}
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(req).Field(i).Interface(), es.Value, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errResponse := utilities.ValidateStruct(tc.params)
			if errResponse != nil {
				for i, es := range *errResponse {
					tc.expected(t, es, tc.params, i)
				}
			}
		})
	}
}

func TestPageRequestValidation(t *testing.T) {
	testCases := []struct {
		name     string
		params   PageRequest
		expected func(t *testing.T, es *utilities.ErrorResponse, req PageRequest, i int)
	}{
		{
			name: "OK",
			params: PageRequest{
				Page: 100,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req PageRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.NotContains(t, es.FailedField, r.Name, m)
				require.NotContains(t, r.Tag, es.Tag, m)
				require.NotEqual(t, reflect.ValueOf(req).Field(i), es.Value, m)
			},
		},
		{
			name:   "Fail_Missing_Fields",
			params: PageRequest{},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req PageRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(req).Field(i).Interface(), es.Value, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errResponse := utilities.ValidateStruct(tc.params)
			if errResponse != nil {
				for i, es := range *errResponse {
					tc.expected(t, es, tc.params, i)
				}
			}
		})
	}
}

func TestMoreInfoRequestValidation(t *testing.T) {
	testCases := []struct {
		name     string
		params   MoreInfoRequest
		expected func(t *testing.T, es *utilities.ErrorResponse, req MoreInfoRequest, i int)
	}{
		{
			name: "OK",
			params: MoreInfoRequest{
				UserId:  "user",
				ScoreId: "123",
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req MoreInfoRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.NotContains(t, es.FailedField, r.Name, m)
				require.NotContains(t, r.Tag, es.Tag, m)
				require.NotEqual(t, reflect.ValueOf(req).Field(i), es.Value, m)
			},
		},
		{
			name: "Fail_Missing_Fields",
			params: MoreInfoRequest{
				ScoreId: "score",
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req MoreInfoRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(req).Field(i).Interface(), es.Value, m)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errResponse := utilities.ValidateStruct(tc.params)
			if errResponse != nil {
				for i, es := range *errResponse {
					tc.expected(t, es, tc.params, i)
				}
			}
		})
	}
}

func TestMetamaskAddrFormat(t *testing.T) {
	r := new(MetamaskAddressRequest)
	type testcase struct {
		addr    string
		success bool
	}
	tcs := []testcase{{"0x1234BF77D1De9eacf66FE81a09a86CfAb212a542", true}, {"0x1234BF77D1De9eacf66FE81a09a86CfAb212a54", false}, {"cx1234BF77D1De9eacf66FE81a09a86CfAb212a542", false}, {"fail", false}}
	for _, tc := range tcs {
		r.Addr = tc.addr
		errs := utilities.ValidateStruct(r)
		if tc.success {
			require.Nil(t, errs)
		} else {
			require.NotNil(t, errs)
		}
	}
}
