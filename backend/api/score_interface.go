package api

import (
	"bitmoi/backend/utilities/common"
	"fmt"
	"math"
)

type ScoreReqInterface interface {
	GetMode() string
	GetUserID() string
	GetScoreID() string
	GetPairName() string
	GetStage() int32
	GetEntryPrice() float64
	GetProfitPrice() float64
	GetLossPrice() float64
	GetIsLong() bool
	GetQuantity() float64
	GetLeverage() int32
	GetBalance() float64
}

type ScoreRequest struct {
	Mode        string  `json:"mode" validate:"required,oneof=competition practice"`
	UserId      string  `json:"user_id" validate:"required,alphanum"`
	ScoreId     string  `json:"score_id" validate:"required,numeric"`
	Name        string  `json:"name" validate:"required"`
	Stage       int32   `json:"stage" validate:"required,number,min=1,max=10"`
	IsLong      *bool   `json:"is_long"  validate:"required,boolean"`
	EntryPrice  float64 `json:"entry_price" validate:"required,number,gt=0"`
	Quantity    float64 `json:"quantity" validate:"required,number,gt=0"`
	ProfitPrice float64 `json:"profit_price" validate:"number,min=0"`
	LossPrice   float64 `json:"loss_price" validate:"number,min=0"`
	Leverage    int32   `json:"leverage" validate:"required,number,min=1,max=100"`
	Balance     float64 `json:"balance" validate:"required,number,gt=0"`
	Identifier  string  `json:"identifier"  validate:"required"`
	WaitingTerm int32   `json:"waiting_term" validate:"required,number,min=1,max=1"`
}

func (s *ScoreRequest) GetMode() string         { return s.Mode }
func (s *ScoreRequest) GetUserID() string       { return s.UserId }
func (s *ScoreRequest) GetScoreID() string      { return s.ScoreId }
func (s *ScoreRequest) GetPairName() string     { return s.Name }
func (s *ScoreRequest) GetStage() int32         { return s.Stage }
func (s *ScoreRequest) GetEntryPrice() float64  { return s.EntryPrice }
func (s *ScoreRequest) GetProfitPrice() float64 { return s.ProfitPrice }
func (s *ScoreRequest) GetLossPrice() float64   { return s.LossPrice }
func (s *ScoreRequest) GetIsLong() bool         { return *s.IsLong }
func (s *ScoreRequest) GetQuantity() float64    { return s.Quantity }
func (s *ScoreRequest) GetLeverage() int32      { return s.Leverage }
func (s *ScoreRequest) GetBalance() float64     { return s.Balance }

type InterScoreRequest struct {
	Mode         string  `json:"mode" validate:"required,oneof=competition practice"`
	UserId       string  `json:"user_id" validate:"required,alphanum"`
	ScoreId      string  `json:"score_id" validate:"required,numeric"`
	Name         string  `json:"name" validate:"required"`
	Stage        int32   `json:"stage" validate:"required,number,min=1,max=10"`
	IsLong       *bool   `json:"is_long"  validate:"required,boolean"`
	EntryPrice   float64 `json:"entry_price" validate:"required,number,gt=0"`
	Quantity     float64 `json:"quantity" validate:"required,number,gt=0"`
	ProfitPrice  float64 `json:"profit_price" validate:"number,min=0"`
	LossPrice    float64 `json:"loss_price" validate:"number,min=0"`
	Leverage     int32   `json:"leverage" validate:"required,number,min=1,max=100"`
	Balance      float64 `json:"balance" validate:"required,number,gt=0"`
	Identifier   string  `json:"identifier"  validate:"required"`
	ReqInterval  string  `json:"reqinterval" validate:"required,oneof=5m 15m 1h 4h 1d" query:"reqinterval"`
	MinTimestamp int64   `json:"min_timestamp" validate:"required,number" query:"min_timestamp"`
	MaxTimestamp int64   `json:"max_timestamp" validate:"required,number" query:"max_timestamp"`
}

func (i *InterScoreRequest) GetMode() string         { return i.Mode }
func (i *InterScoreRequest) GetUserID() string       { return i.UserId }
func (i *InterScoreRequest) GetScoreID() string      { return i.ScoreId }
func (i *InterScoreRequest) GetPairName() string     { return i.Name }
func (i *InterScoreRequest) GetStage() int32         { return i.Stage }
func (i *InterScoreRequest) GetEntryPrice() float64  { return i.EntryPrice }
func (i *InterScoreRequest) GetProfitPrice() float64 { return i.ProfitPrice }
func (i *InterScoreRequest) GetLossPrice() float64   { return i.LossPrice }
func (i *InterScoreRequest) GetIsLong() bool         { return *i.IsLong }
func (i *InterScoreRequest) GetQuantity() float64    { return i.Quantity }
func (i *InterScoreRequest) GetLeverage() int32      { return i.Leverage }
func (i *InterScoreRequest) GetBalance() float64     { return i.Balance }

func validateOrderRequest(s ScoreReqInterface) error {
	var (
		entryPrice  = s.GetEntryPrice()
		profitPrice = s.GetProfitPrice()
		lossPrice   = s.GetLossPrice()
		leverage    = s.GetLeverage()
		quantity    = s.GetQuantity()
		balance     = s.GetBalance()
	)

	limit := math.Pow(float64(leverage), float64(-1))
	if s.GetIsLong() {
		if !(entryPrice < profitPrice && lossPrice < entryPrice) {
			return fmt.Errorf("check the profit and loss price. positon: %s, entry price: %f", long, entryPrice)
		}
		if (entryPrice-lossPrice)/entryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				long, entryPrice, lossPrice, leverage, common.CeilDecimal(entryPrice*(1-limit)))
		}
	} else {
		if !(entryPrice > profitPrice && lossPrice > entryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : %s, entry price : %f", short, entryPrice)
		}
		if (lossPrice-entryPrice)/entryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				short, entryPrice, lossPrice, leverage, common.FloorDecimal(entryPrice*(1+limit)))
		}
	}

	if (balance * float64(leverage)) < (quantity * entryPrice) {
		return fmt.Errorf("invalid order. check your balance. order amount : %.5f, limit amount : %.5f ", quantity*entryPrice, balance*float64(leverage))
	}
	return nil
}

type OrderResultInterface interface {
	GetPairName() string
	GetEntryTime() string
	GetEndPrice() float64
	GetOutTime() int64
	GetRoe() float64
	GetPnl() float64
	GetIsliquidated() bool
}

type OrderResult struct {
	Name         string  `json:"name"`
	Entrytime    string  `json:"entry_time"`
	Leverage     int32   `json:"leverage"`
	EndPrice     float64 `json:"end_price"`
	OutTime      int32   `json:"out_time"`
	Roe          float64 `json:"roe"`
	Pnl          float64 `json:"pnl"`
	Commission   float64 `json:"commission"`
	Isliquidated bool    `json:"is_liquidated"`
}

func (o *OrderResult) GetPairName() string   { return o.Name }
func (o *OrderResult) GetEntryTime() string  { return o.Entrytime }
func (o *OrderResult) GetEndPrice() float64  { return o.EndPrice }
func (o *OrderResult) GetOutTime() int64     { return int64(o.OutTime) }
func (o *OrderResult) GetRoe() float64       { return o.Roe }
func (o *OrderResult) GetPnl() float64       { return o.Pnl }
func (o *OrderResult) GetIsliquidated() bool { return o.Isliquidated }

type InterMediateResult struct {
	Name         string  `json:"name"`
	Entrytime    string  `json:"entry_time"`
	Leverage     int32   `json:"leverage"`
	EndPrice     float64 `json:"end_price"`
	OutTime      int64   `json:"out_time"`
	Roe          float64 `json:"roe"`
	Pnl          float64 `json:"pnl"`
	Commission   float64 `json:"commission"`
	Isliquidated bool    `json:"is_liquidated"`
}

func (i *InterMediateResult) GetPairName() string   { return i.Name }
func (i *InterMediateResult) GetEntryTime() string  { return i.Entrytime }
func (i *InterMediateResult) GetEndPrice() float64  { return i.EndPrice }
func (i *InterMediateResult) GetOutTime() int64     { return i.OutTime }
func (i *InterMediateResult) GetRoe() float64       { return i.Roe }
func (i *InterMediateResult) GetPnl() float64       { return i.Pnl }
func (i *InterMediateResult) GetIsliquidated() bool { return i.Isliquidated }
