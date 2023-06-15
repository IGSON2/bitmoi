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
		params   OrderRequest
		expected func(t *testing.T, es *utilities.ErrorResponse, o OrderRequest, i int)
	}{
		{
			name: "OK",
			params: OrderRequest{
				Mode:         competition,
				UserId:       "user",
				Name:         "user_name",
				Stage:        1,
				IsLong:       &short,
				EntryPrice:   100.123,
				Quantity:     0.123,
				QuantityRate: 15.75,
				ProfitPrice:  110.2521,
				LossPrice:    90.2817,
				Leverage:     25,
				Balance:      1000,
				Identifier:   "asdf",
				ScoreId:      "12345",
				WaitingTerm:  1,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o OrderRequest, i int) {
				require.Nil(t, es)
			},
		},
		{
			name: "Fail_Missing_Fields",
			params: OrderRequest{
				Mode:         "",
				UserId:       "",
				Name:         "",
				Stage:        -1,
				IsLong:       nil,
				EntryPrice:   -1,
				Quantity:     -1,
				QuantityRate: -1,
				ProfitPrice:  -1,
				LossPrice:    -1,
				Leverage:     0,
				Balance:      -1,
				Identifier:   "",
				ScoreId:      "",
				WaitingTerm:  0,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o OrderRequest, i int) {
				r := reflect.TypeOf(o).Field(i)
				m := fmt.Sprintf("Field : %s, Tag : %s, Value : %s", es.FailedField, es.Tag, es.Value)
				require.Contains(t, es.FailedField, r.Name, m)
				require.Contains(t, r.Tag, es.Tag, m)
				require.Equal(t, reflect.ValueOf(o).Field(i).Interface(), es.Value, m)
			},
		},
		{
			name: "Fail_Boundary_Value",
			params: OrderRequest{
				Mode:         "",
				UserId:       "",
				Name:         "",
				Stage:        0,
				IsLong:       nil,
				EntryPrice:   0,
				Quantity:     0,
				QuantityRate: 0,
				ProfitPrice:  0,
				LossPrice:    0,
				Leverage:     0,
				Balance:      0,
				Identifier:   "",
				ScoreId:      "",
				WaitingTerm:  0,
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, o OrderRequest, i int) {
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
				UserId:      "user",
				ScoreId:     "123",
				Comment:     "comment",
				DisplayName: "name",
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
				UserId:      "",
				ScoreId:     "score",
				Comment:     "",
				DisplayName: "name",
			},
			expected: func(t *testing.T, es *utilities.ErrorResponse, req RankInsertRequest, i int) {
				r := reflect.TypeOf(req).Field(i)
				if r.Name == "Comment" || r.Name == "DisplayName" {
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
