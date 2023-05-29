// Code generated by MockGen. DO NOT EDIT.
// Source: bitmoi/backend/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	db "bitmoi/backend/db/sqlc"
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), arg0, arg1)
}

// Get15mCandles mocks base method.
func (m *MockStore) Get15mCandles(arg0 context.Context, arg1 db.Get15mCandlesParams) ([]db.Candles15m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get15mCandles", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles15m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get15mCandles indicates an expected call of Get15mCandles.
func (mr *MockStoreMockRecorder) Get15mCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get15mCandles", reflect.TypeOf((*MockStore)(nil).Get15mCandles), arg0, arg1)
}

// Get15mMinMaxTime mocks base method.
func (m *MockStore) Get15mMinMaxTime(arg0 context.Context, arg1 string) (db.Get15mMinMaxTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get15mMinMaxTime", arg0, arg1)
	ret0, _ := ret[0].(db.Get15mMinMaxTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get15mMinMaxTime indicates an expected call of Get15mMinMaxTime.
func (mr *MockStoreMockRecorder) Get15mMinMaxTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get15mMinMaxTime", reflect.TypeOf((*MockStore)(nil).Get15mMinMaxTime), arg0, arg1)
}

// Get15mResult mocks base method.
func (m *MockStore) Get15mResult(arg0 context.Context, arg1 db.Get15mResultParams) ([]db.Candles15m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get15mResult", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles15m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get15mResult indicates an expected call of Get15mResult.
func (mr *MockStoreMockRecorder) Get15mResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get15mResult", reflect.TypeOf((*MockStore)(nil).Get15mResult), arg0, arg1)
}

// Get15mVolSumPriceAVG mocks base method.
func (m *MockStore) Get15mVolSumPriceAVG(arg0 context.Context, arg1 db.Get15mVolSumPriceAVGParams) (db.Get15mVolSumPriceAVGRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get15mVolSumPriceAVG", arg0, arg1)
	ret0, _ := ret[0].(db.Get15mVolSumPriceAVGRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get15mVolSumPriceAVG indicates an expected call of Get15mVolSumPriceAVG.
func (mr *MockStoreMockRecorder) Get15mVolSumPriceAVG(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get15mVolSumPriceAVG", reflect.TypeOf((*MockStore)(nil).Get15mVolSumPriceAVG), arg0, arg1)
}

// Get1dCandles mocks base method.
func (m *MockStore) Get1dCandles(arg0 context.Context, arg1 db.Get1dCandlesParams) ([]db.Candles1d, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1dCandles", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles1d)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1dCandles indicates an expected call of Get1dCandles.
func (mr *MockStoreMockRecorder) Get1dCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1dCandles", reflect.TypeOf((*MockStore)(nil).Get1dCandles), arg0, arg1)
}

// Get1dMinMaxTime mocks base method.
func (m *MockStore) Get1dMinMaxTime(arg0 context.Context, arg1 string) (db.Get1dMinMaxTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1dMinMaxTime", arg0, arg1)
	ret0, _ := ret[0].(db.Get1dMinMaxTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1dMinMaxTime indicates an expected call of Get1dMinMaxTime.
func (mr *MockStoreMockRecorder) Get1dMinMaxTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1dMinMaxTime", reflect.TypeOf((*MockStore)(nil).Get1dMinMaxTime), arg0, arg1)
}

// Get1dResult mocks base method.
func (m *MockStore) Get1dResult(arg0 context.Context, arg1 db.Get1dResultParams) ([]db.Candles1d, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1dResult", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles1d)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1dResult indicates an expected call of Get1dResult.
func (mr *MockStoreMockRecorder) Get1dResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1dResult", reflect.TypeOf((*MockStore)(nil).Get1dResult), arg0, arg1)
}

// Get1dVolSumPriceAVG mocks base method.
func (m *MockStore) Get1dVolSumPriceAVG(arg0 context.Context, arg1 db.Get1dVolSumPriceAVGParams) (db.Get1dVolSumPriceAVGRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1dVolSumPriceAVG", arg0, arg1)
	ret0, _ := ret[0].(db.Get1dVolSumPriceAVGRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1dVolSumPriceAVG indicates an expected call of Get1dVolSumPriceAVG.
func (mr *MockStoreMockRecorder) Get1dVolSumPriceAVG(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1dVolSumPriceAVG", reflect.TypeOf((*MockStore)(nil).Get1dVolSumPriceAVG), arg0, arg1)
}

// Get1hCandles mocks base method.
func (m *MockStore) Get1hCandles(arg0 context.Context, arg1 db.Get1hCandlesParams) ([]db.Candles1h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1hCandles", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles1h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1hCandles indicates an expected call of Get1hCandles.
func (mr *MockStoreMockRecorder) Get1hCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1hCandles", reflect.TypeOf((*MockStore)(nil).Get1hCandles), arg0, arg1)
}

// Get1hMinMaxTime mocks base method.
func (m *MockStore) Get1hMinMaxTime(arg0 context.Context, arg1 string) (db.Get1hMinMaxTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1hMinMaxTime", arg0, arg1)
	ret0, _ := ret[0].(db.Get1hMinMaxTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1hMinMaxTime indicates an expected call of Get1hMinMaxTime.
func (mr *MockStoreMockRecorder) Get1hMinMaxTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1hMinMaxTime", reflect.TypeOf((*MockStore)(nil).Get1hMinMaxTime), arg0, arg1)
}

// Get1hResult mocks base method.
func (m *MockStore) Get1hResult(arg0 context.Context, arg1 db.Get1hResultParams) ([]db.Candles1h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1hResult", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles1h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1hResult indicates an expected call of Get1hResult.
func (mr *MockStoreMockRecorder) Get1hResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1hResult", reflect.TypeOf((*MockStore)(nil).Get1hResult), arg0, arg1)
}

// Get1hVolSumPriceAVG mocks base method.
func (m *MockStore) Get1hVolSumPriceAVG(arg0 context.Context, arg1 db.Get1hVolSumPriceAVGParams) (db.Get1hVolSumPriceAVGRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get1hVolSumPriceAVG", arg0, arg1)
	ret0, _ := ret[0].(db.Get1hVolSumPriceAVGRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get1hVolSumPriceAVG indicates an expected call of Get1hVolSumPriceAVG.
func (mr *MockStoreMockRecorder) Get1hVolSumPriceAVG(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get1hVolSumPriceAVG", reflect.TypeOf((*MockStore)(nil).Get1hVolSumPriceAVG), arg0, arg1)
}

// Get4hCandles mocks base method.
func (m *MockStore) Get4hCandles(arg0 context.Context, arg1 db.Get4hCandlesParams) ([]db.Candles4h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get4hCandles", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles4h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get4hCandles indicates an expected call of Get4hCandles.
func (mr *MockStoreMockRecorder) Get4hCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get4hCandles", reflect.TypeOf((*MockStore)(nil).Get4hCandles), arg0, arg1)
}

// Get4hMinMaxTime mocks base method.
func (m *MockStore) Get4hMinMaxTime(arg0 context.Context, arg1 string) (db.Get4hMinMaxTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get4hMinMaxTime", arg0, arg1)
	ret0, _ := ret[0].(db.Get4hMinMaxTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get4hMinMaxTime indicates an expected call of Get4hMinMaxTime.
func (mr *MockStoreMockRecorder) Get4hMinMaxTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get4hMinMaxTime", reflect.TypeOf((*MockStore)(nil).Get4hMinMaxTime), arg0, arg1)
}

// Get4hResult mocks base method.
func (m *MockStore) Get4hResult(arg0 context.Context, arg1 db.Get4hResultParams) ([]db.Candles4h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get4hResult", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles4h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get4hResult indicates an expected call of Get4hResult.
func (mr *MockStoreMockRecorder) Get4hResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get4hResult", reflect.TypeOf((*MockStore)(nil).Get4hResult), arg0, arg1)
}

// Get4hVolSumPriceAVG mocks base method.
func (m *MockStore) Get4hVolSumPriceAVG(arg0 context.Context, arg1 db.Get4hVolSumPriceAVGParams) (db.Get4hVolSumPriceAVGRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get4hVolSumPriceAVG", arg0, arg1)
	ret0, _ := ret[0].(db.Get4hVolSumPriceAVGRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get4hVolSumPriceAVG indicates an expected call of Get4hVolSumPriceAVG.
func (mr *MockStoreMockRecorder) Get4hVolSumPriceAVG(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get4hVolSumPriceAVG", reflect.TypeOf((*MockStore)(nil).Get4hVolSumPriceAVG), arg0, arg1)
}

// Get5mCandles mocks base method.
func (m *MockStore) Get5mCandles(arg0 context.Context, arg1 db.Get5mCandlesParams) ([]db.Candles5m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get5mCandles", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles5m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get5mCandles indicates an expected call of Get5mCandles.
func (mr *MockStoreMockRecorder) Get5mCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get5mCandles", reflect.TypeOf((*MockStore)(nil).Get5mCandles), arg0, arg1)
}

// Get5mMinMaxTime mocks base method.
func (m *MockStore) Get5mMinMaxTime(arg0 context.Context, arg1 string) (db.Get5mMinMaxTimeRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get5mMinMaxTime", arg0, arg1)
	ret0, _ := ret[0].(db.Get5mMinMaxTimeRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get5mMinMaxTime indicates an expected call of Get5mMinMaxTime.
func (mr *MockStoreMockRecorder) Get5mMinMaxTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get5mMinMaxTime", reflect.TypeOf((*MockStore)(nil).Get5mMinMaxTime), arg0, arg1)
}

// Get5mResult mocks base method.
func (m *MockStore) Get5mResult(arg0 context.Context, arg1 db.Get5mResultParams) ([]db.Candles5m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get5mResult", arg0, arg1)
	ret0, _ := ret[0].([]db.Candles5m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get5mResult indicates an expected call of Get5mResult.
func (mr *MockStoreMockRecorder) Get5mResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get5mResult", reflect.TypeOf((*MockStore)(nil).Get5mResult), arg0, arg1)
}

// Get5mVolSumPriceAVG mocks base method.
func (m *MockStore) Get5mVolSumPriceAVG(arg0 context.Context, arg1 db.Get5mVolSumPriceAVGParams) (db.Get5mVolSumPriceAVGRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get5mVolSumPriceAVG", arg0, arg1)
	ret0, _ := ret[0].(db.Get5mVolSumPriceAVGRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get5mVolSumPriceAVG indicates an expected call of Get5mVolSumPriceAVG.
func (mr *MockStoreMockRecorder) Get5mVolSumPriceAVG(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get5mVolSumPriceAVG", reflect.TypeOf((*MockStore)(nil).Get5mVolSumPriceAVG), arg0, arg1)
}

// GetAllParisInDB mocks base method.
func (m *MockStore) GetAllParisInDB(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllParisInDB", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllParisInDB indicates an expected call of GetAllParisInDB.
func (mr *MockStoreMockRecorder) GetAllParisInDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllParisInDB", reflect.TypeOf((*MockStore)(nil).GetAllParisInDB), arg0)
}

// GetAllRanks mocks base method.
func (m *MockStore) GetAllRanks(arg0 context.Context, arg1 db.GetAllRanksParams) ([]db.RankingBoard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRanks", arg0, arg1)
	ret0, _ := ret[0].([]db.RankingBoard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllRanks indicates an expected call of GetAllRanks.
func (mr *MockStoreMockRecorder) GetAllRanks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRanks", reflect.TypeOf((*MockStore)(nil).GetAllRanks), arg0, arg1)
}

// GetLastUser mocks base method.
func (m *MockStore) GetLastUser(arg0 context.Context) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastUser", arg0)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastUser indicates an expected call of GetLastUser.
func (mr *MockStoreMockRecorder) GetLastUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastUser", reflect.TypeOf((*MockStore)(nil).GetLastUser), arg0)
}

// GetOne15mCandle mocks base method.
func (m *MockStore) GetOne15mCandle(arg0 context.Context, arg1 db.GetOne15mCandleParams) (db.Candles15m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne15mCandle", arg0, arg1)
	ret0, _ := ret[0].(db.Candles15m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne15mCandle indicates an expected call of GetOne15mCandle.
func (mr *MockStoreMockRecorder) GetOne15mCandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne15mCandle", reflect.TypeOf((*MockStore)(nil).GetOne15mCandle), arg0, arg1)
}

// GetOne1dCandle mocks base method.
func (m *MockStore) GetOne1dCandle(arg0 context.Context, arg1 db.GetOne1dCandleParams) (db.Candles1d, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne1dCandle", arg0, arg1)
	ret0, _ := ret[0].(db.Candles1d)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne1dCandle indicates an expected call of GetOne1dCandle.
func (mr *MockStoreMockRecorder) GetOne1dCandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne1dCandle", reflect.TypeOf((*MockStore)(nil).GetOne1dCandle), arg0, arg1)
}

// GetOne1hCandle mocks base method.
func (m *MockStore) GetOne1hCandle(arg0 context.Context, arg1 db.GetOne1hCandleParams) (db.Candles1h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne1hCandle", arg0, arg1)
	ret0, _ := ret[0].(db.Candles1h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne1hCandle indicates an expected call of GetOne1hCandle.
func (mr *MockStoreMockRecorder) GetOne1hCandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne1hCandle", reflect.TypeOf((*MockStore)(nil).GetOne1hCandle), arg0, arg1)
}

// GetOne4hCandle mocks base method.
func (m *MockStore) GetOne4hCandle(arg0 context.Context, arg1 db.GetOne4hCandleParams) (db.Candles4h, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne4hCandle", arg0, arg1)
	ret0, _ := ret[0].(db.Candles4h)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne4hCandle indicates an expected call of GetOne4hCandle.
func (mr *MockStoreMockRecorder) GetOne4hCandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne4hCandle", reflect.TypeOf((*MockStore)(nil).GetOne4hCandle), arg0, arg1)
}

// GetOne5mCandle mocks base method.
func (m *MockStore) GetOne5mCandle(arg0 context.Context, arg1 db.GetOne5mCandleParams) (db.Candles5m, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne5mCandle", arg0, arg1)
	ret0, _ := ret[0].(db.Candles5m)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne5mCandle indicates an expected call of GetOne5mCandle.
func (mr *MockStoreMockRecorder) GetOne5mCandle(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne5mCandle", reflect.TypeOf((*MockStore)(nil).GetOne5mCandle), arg0, arg1)
}

// GetRandomUser mocks base method.
func (m *MockStore) GetRandomUser(arg0 context.Context) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRandomUser", arg0)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRandomUser indicates an expected call of GetRandomUser.
func (mr *MockStoreMockRecorder) GetRandomUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRandomUser", reflect.TypeOf((*MockStore)(nil).GetRandomUser), arg0)
}

// GetRankByUserID mocks base method.
func (m *MockStore) GetRankByUserID(arg0 context.Context, arg1 string) (db.RankingBoard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRankByUserID", arg0, arg1)
	ret0, _ := ret[0].(db.RankingBoard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRankByUserID indicates an expected call of GetRankByUserID.
func (mr *MockStoreMockRecorder) GetRankByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRankByUserID", reflect.TypeOf((*MockStore)(nil).GetRankByUserID), arg0, arg1)
}

// GetScore mocks base method.
func (m *MockStore) GetScore(arg0 context.Context, arg1 db.GetScoreParams) (db.Score, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScore", arg0, arg1)
	ret0, _ := ret[0].(db.Score)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScore indicates an expected call of GetScore.
func (mr *MockStoreMockRecorder) GetScore(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScore", reflect.TypeOf((*MockStore)(nil).GetScore), arg0, arg1)
}

// GetScoreToStage mocks base method.
func (m *MockStore) GetScoreToStage(arg0 context.Context, arg1 db.GetScoreToStageParams) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScoreToStage", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScoreToStage indicates an expected call of GetScoreToStage.
func (mr *MockStoreMockRecorder) GetScoreToStage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScoreToStage", reflect.TypeOf((*MockStore)(nil).GetScoreToStage), arg0, arg1)
}

// GetScoresByScoreID mocks base method.
func (m *MockStore) GetScoresByScoreID(arg0 context.Context, arg1 db.GetScoresByScoreIDParams) ([]db.Score, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScoresByScoreID", arg0, arg1)
	ret0, _ := ret[0].([]db.Score)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScoresByScoreID indicates an expected call of GetScoresByScoreID.
func (mr *MockStoreMockRecorder) GetScoresByScoreID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScoresByScoreID", reflect.TypeOf((*MockStore)(nil).GetScoresByScoreID), arg0, arg1)
}

// GetScoresByUserID mocks base method.
func (m *MockStore) GetScoresByUserID(arg0 context.Context, arg1 db.GetScoresByUserIDParams) ([]db.Score, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetScoresByUserID", arg0, arg1)
	ret0, _ := ret[0].([]db.Score)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetScoresByUserID indicates an expected call of GetScoresByUserID.
func (mr *MockStoreMockRecorder) GetScoresByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetScoresByUserID", reflect.TypeOf((*MockStore)(nil).GetScoresByUserID), arg0, arg1)
}

// GetStageLenByScoreID mocks base method.
func (m *MockStore) GetStageLenByScoreID(arg0 context.Context, arg1 db.GetStageLenByScoreIDParams) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStageLenByScoreID", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStageLenByScoreID indicates an expected call of GetStageLenByScoreID.
func (mr *MockStoreMockRecorder) GetStageLenByScoreID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStageLenByScoreID", reflect.TypeOf((*MockStore)(nil).GetStageLenByScoreID), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), arg0, arg1)
}

// Insert15mCandles mocks base method.
func (m *MockStore) Insert15mCandles(arg0 context.Context, arg1 db.Insert15mCandlesParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert15mCandles", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert15mCandles indicates an expected call of Insert15mCandles.
func (mr *MockStoreMockRecorder) Insert15mCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert15mCandles", reflect.TypeOf((*MockStore)(nil).Insert15mCandles), arg0, arg1)
}

// Insert1dCandles mocks base method.
func (m *MockStore) Insert1dCandles(arg0 context.Context, arg1 db.Insert1dCandlesParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert1dCandles", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert1dCandles indicates an expected call of Insert1dCandles.
func (mr *MockStoreMockRecorder) Insert1dCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert1dCandles", reflect.TypeOf((*MockStore)(nil).Insert1dCandles), arg0, arg1)
}

// Insert1hCandles mocks base method.
func (m *MockStore) Insert1hCandles(arg0 context.Context, arg1 db.Insert1hCandlesParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert1hCandles", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert1hCandles indicates an expected call of Insert1hCandles.
func (mr *MockStoreMockRecorder) Insert1hCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert1hCandles", reflect.TypeOf((*MockStore)(nil).Insert1hCandles), arg0, arg1)
}

// Insert4hCandles mocks base method.
func (m *MockStore) Insert4hCandles(arg0 context.Context, arg1 db.Insert4hCandlesParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert4hCandles", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert4hCandles indicates an expected call of Insert4hCandles.
func (mr *MockStoreMockRecorder) Insert4hCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert4hCandles", reflect.TypeOf((*MockStore)(nil).Insert4hCandles), arg0, arg1)
}

// Insert5mCandles mocks base method.
func (m *MockStore) Insert5mCandles(arg0 context.Context, arg1 db.Insert5mCandlesParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert5mCandles", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert5mCandles indicates an expected call of Insert5mCandles.
func (mr *MockStoreMockRecorder) Insert5mCandles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert5mCandles", reflect.TypeOf((*MockStore)(nil).Insert5mCandles), arg0, arg1)
}

// InsertRank mocks base method.
func (m *MockStore) InsertRank(arg0 context.Context, arg1 db.InsertRankParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertRank", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertRank indicates an expected call of InsertRank.
func (mr *MockStoreMockRecorder) InsertRank(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertRank", reflect.TypeOf((*MockStore)(nil).InsertRank), arg0, arg1)
}

// InsertScore mocks base method.
func (m *MockStore) InsertScore(arg0 context.Context, arg1 db.InsertScoreParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertScore", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertScore indicates an expected call of InsertScore.
func (mr *MockStoreMockRecorder) InsertScore(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertScore", reflect.TypeOf((*MockStore)(nil).InsertScore), arg0, arg1)
}

// UpdatePhotoURL mocks base method.
func (m *MockStore) UpdatePhotoURL(arg0 context.Context, arg1 db.UpdatePhotoURLParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePhotoURL", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePhotoURL indicates an expected call of UpdatePhotoURL.
func (mr *MockStoreMockRecorder) UpdatePhotoURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePhotoURL", reflect.TypeOf((*MockStore)(nil).UpdatePhotoURL), arg0, arg1)
}

// UpdateUserRank mocks base method.
func (m *MockStore) UpdateUserRank(arg0 context.Context, arg1 db.UpdateUserRankParams) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUserRank", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUserRank indicates an expected call of UpdateUserRank.
func (mr *MockStoreMockRecorder) UpdateUserRank(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUserRank", reflect.TypeOf((*MockStore)(nil).UpdateUserRank), arg0, arg1)
}
