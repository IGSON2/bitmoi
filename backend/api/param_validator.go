package api

type CandlesRequest struct {
	Names string `json:"names" query:"names"`
}

type RankInsertRequest struct {
	ScoreId string `json:"score_id" validate:"required,numeric"`
	Comment string `json:"comment"`
}

type MyscoreRequest struct {
	Mode string `json:"mode" validate:"required,oneof=competition practice" query:"mode"`
	Page int32  `json:"page" validate:"min=1,number" query:"page"`
}

type MoreInfoRequest struct {
	UserId  string `json:"user_id" validate:"required,alphanum" query:"userid"`
	ScoreId string `json:"score_id" validate:"required,numeric" query:"scoreid"`
}

type AnotherIntervalRequest struct {
	ReqInterval string `json:"reqinterval" validate:"required,oneof=5m 15m 1h 4h 1d" query:"reqinterval"`
	Identifier  string `json:"identifier" validate:"required" query:"identifier"`
	Mode        string `json:"mode" validate:"required,oneof=competition practice" query:"mode"`
}

type LoginUserRequest struct {
	UserID   string `json:"user_id" validate:"required,alphanum,min=5,max=15"`
	Password string `json:"password" validate:"required,min=8"`
}

type CreateUserRequest struct {
	UserID   string `json:"user_id" validate:"required,alphanum,min=5,max=15"`
	Password string `json:"password" validate:"required,min=8"`
	Nickname string `json:"nickname" validate:"required,min=1,max=10"`
	Email    string `json:"email" validate:"required,email"`
	OauthUid string `json:"oauth_uid,omitempty"`
}

type VerifyEmailRequest struct {
	EmailId    int64  `json:"email_id" validate:"required,min=1" query:"email_id"`
	SecretCode string `json:"secret_code" validate:"required,min=32,max=128" query:"secret_code"`
}

type UpdateUsingTokenRequest struct {
	ScoreId string `json:"score_id" validate:"required,numeric"`
}

type UpdateMetamaskRequest struct {
	Addr string `json:"addr" validate:"required,eth_addr"`
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname" validate:"required,min=1,max=10"`
}

type GetBidderByLocRequest struct {
	Location string `json:"location" validate:"required,oneof=practice rank freetoken" query:"location"`
}

type BidTokenRequest struct {
	Amount   int    `json:"amount" validate:"required,number,min=1"`
	Location string `json:"location" validate:"required,oneof=practice rank freetoken"`
}

type CreateRecommendHistoryRequest struct {
	Code string `json:"code" validate:"required,hexadecimal"`
}

type GetRankRequest struct {
	Mode     string `json:"mode" validate:"oneof=competition practice" query:"mode"`
	Category string `json:"category" validate:"oneof=pnl roe" query:"category"`
	Start    string `json:"start" validate:"datetime=06-01-02" query:"start"`
	End      string `json:"end" validate:"datetime=06-01-02" query:"end"`
	Page     int32  `json:"page" validate:"min=1,number" query:"page"`
}
