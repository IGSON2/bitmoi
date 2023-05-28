package api

const (
	competition = "competition"
	practice    = "practice"
)

type ChartRequestQuery struct {
	Names    string `json:"names" validate:"required"`
	Interval string `json:"interval" validate:"required,oneof=5m 15m 1h 4h 1d"`
}

type OrderRequest struct {
	Mode         string  `json:"mode" validate:"required,oneof=competition practice"`
	UserId       string  `json:"userid" validate:"required,alpha"`
	Name         string  `json:"name" validate:"required"`
	Entrytime    string  `json:"entrytime" validate:"required"`
	Stage        int32   `json:"stage" validate:"required,number,min=1,max=10"`
	IsLong       *bool   `json:"islong"  validate:"required,boolean"`
	EntryPrice   float64 `json:"entryprice" validate:"required,number,min=0"`
	Quantity     float64 `json:"quantity" validate:"required,number,min=0"`
	QuantityRate float64 `json:"quantityrate" validate:"required,number,min=0,max=100"`
	ProfitPrice  float64 `json:"profitprice" validate:"required,number,min=0"`
	LossPrice    float64 `json:"lossprice" validate:"required,number,min=0"`
	Leverage     int32   `json:"leverage" validate:"required,number,min=0,max=100"`
	Balance      float64 `json:"balance" validate:"required,number,min=0"`
	Identifier   string  `json:"identifier,omitempty"  validate:"required"`
	ScoreId      string  `json:"scoreid" validate:"required,numeric"`
	WaitingTerm  int32   `json:"waitingterm" validate:"required,number,min=1,max=30"`
}

type RankInsertRequest struct {
	UserId      string `json:"userid" validate:"required,alpha"`
	ScoreId     string `json:"scoreid" validate:"required,numeric"`
	Comment     string `json:"comment"`
	DisplayName string `json:"displayname"`
}

type PageRequest struct {
	Page int32 `json:"page" validate:"required,min=1,number"`
}

type MoreInfoRequest struct {
	UserId  string `json:"userid" validate:"required,alpha"`
	ScoreId string `json:"scoreid" validate:"required,numeric"`
}
