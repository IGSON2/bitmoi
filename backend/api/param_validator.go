package api

import (
	"bitmoi/backend/utilities/common"
	"fmt"
	"math"
)

type CandlesRequest struct {
	Names string `json:"names" query:"names"`
}

type ScoreRequest struct {
	Mode        string  `json:"mode" validate:"required,oneof=competition practice"`
	UserId      string  `json:"user_id" validate:"required,alphanum"`
	Name        string  `json:"name" validate:"required"`
	Stage       int32   `json:"stage" validate:"required,number,min=1,max=10"`
	IsLong      *bool   `json:"is_long"  validate:"required,boolean"`
	EntryPrice  float64 `json:"entry_price" validate:"required,number,gt=0"`
	Quantity    float64 `json:"quantity" validate:"required,number,gt=0"`
	ProfitPrice float64 `json:"profit_price" validate:"required,number,min=0"`
	LossPrice   float64 `json:"loss_price" validate:"required,number,min=0"`
	Leverage    int32   `json:"leverage" validate:"required,number,min=1,max=100"`
	Balance     float64 `json:"balance" validate:"required,number,gt=0"`
	Identifier  string  `json:"identifier"  validate:"required"`
	ScoreId     string  `json:"score_id" validate:"required,numeric"`
	WaitingTerm int32   `json:"waiting_term" validate:"required,number,min=1,max=30"`
}

func validateOrderRequest(o *ScoreRequest) error {
	limit := math.Pow(float64(o.Leverage), float64(-1))
	if *o.IsLong {
		if !(o.EntryPrice < o.ProfitPrice && o.LossPrice < o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon: %s, entry price: %f", long, o.EntryPrice)
		}
		if (o.EntryPrice-o.LossPrice)/o.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				long, o.EntryPrice, o.LossPrice, o.Leverage, common.CeilDecimal(o.EntryPrice*(1-limit)))
		}
	} else {
		if !(o.EntryPrice > o.ProfitPrice && o.LossPrice > o.EntryPrice) {
			return fmt.Errorf("check the profit and loss price. positon : %s, entry price : %f", short, o.EntryPrice)
		}
		if (o.LossPrice-o.EntryPrice)/o.EntryPrice > limit {
			return fmt.Errorf("unacceptable loss price. position: %s, entry price: %f, loss price: %f, leverage: %d limit : %f",
				short, o.EntryPrice, o.LossPrice, o.Leverage, common.FloorDecimal(o.EntryPrice*(1+limit)))
		}
	}

	if (o.Balance * float64(o.Leverage)) < (o.Quantity * o.EntryPrice) {
		return fmt.Errorf("invalid order. check your balance. order amount : %.5f, limit amount : %.5f ", o.Quantity*o.EntryPrice, o.Balance*float64(o.Leverage))
	}
	return nil
}

type RankInsertRequest struct {
	UserId      string `json:"userid" validate:"required,alpha"`
	ScoreId     string `json:"scoreid" validate:"required,numeric"`
	Comment     string `json:"comment"`
	DisplayName string `json:"displayname"`
}

type PageRequest struct {
	Page int32 `json:"page" validate:"min=0,number" query:"page"`
}

type MoreInfoRequest struct {
	UserId  string `json:"user_id" validate:"required,alpha" query:"userid"`
	ScoreId string `json:"score_id" validate:"required,numeric" query:"scoreid"`
}

type AnotherIntervalRequest struct {
	ReqInterval string `json:"reqinterval" validate:"required,oneof=5m 15m 1h 4h 1d" query:"reqinterval"`
	Identifier  string `json:"identifier" validate:"required" query:"identifier"`
	Mode        string `json:"mode" validate:"required,oneof=competition practice" query:"mode"`
	Stage       int32  `json:"stage" validate:"required,number,min=1,max=10" query:"stage"`
}

type LoginUserRequest struct {
	UserID   string `json:"user_id" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreateUserRequest struct {
	UserID   string `json:"user_id" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	PhotoUrl string `json:"photo_url,omitempty"`
	OauthUid string `json:"oauth_uid,omitempty"`
}

type VerifyEmailRequest struct {
	EmailId    int64  `json:"email_id" validate:"required,min=1"`
	SecretCode string `json:"secret_code" validate:"required,min=32,max=128"`
}
