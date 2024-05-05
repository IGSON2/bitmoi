// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"
)

type AccumulationHistory struct {
	ToUser    string    `json:"to_user"`
	Amount    float64   `json:"amount"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"method"`
	Giver     string    `json:"giver"`
}

type BiddingHistory struct {
	TxHash    string    `json:"tx_hash"`
	UserID    string    `json:"user_id"`
	Amount    int64     `json:"amount"`
	Location  string    `json:"location"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Candles15m struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

type Candles1d struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

type Candles1h struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

type Candles4h struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

type Candles5m struct {
	Name   string  `json:"name"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Time   int64   `json:"time"`
	Volume float64 `json:"volume"`
	Color  string  `json:"color"`
}

type CompScore struct {
	ScoreID    string       `json:"score_id"`
	UserID     string       `json:"user_id"`
	Stage      int8         `json:"stage"`
	Pairname   string       `json:"pairname"`
	Entrytime  string       `json:"entrytime"`
	Position   string       `json:"position"`
	Leverage   int8         `json:"leverage"`
	Outtime    string       `json:"outtime"`
	Entryprice float64      `json:"entryprice"`
	Quantity   float64      `json:"quantity"`
	Endprice   float64      `json:"endprice"`
	Pnl        float64      `json:"pnl"`
	Roe        float64      `json:"roe"`
	SettledAt  sql.NullTime `json:"settled_at"`
	CreatedAt  time.Time    `json:"created_at"`
}

type PracAfterScore struct {
	ScoreID      string  `json:"score_id"`
	UserID       string  `json:"user_id"`
	MaxRoe       float64 `json:"max_roe"`
	MinRoe       float64 `json:"min_roe"`
	AfterOuttime int64   `json:"after_outtime"`
}

type PracScore struct {
	ScoreID    string       `json:"score_id"`
	UserID     string       `json:"user_id"`
	Stage      int8         `json:"stage"`
	Pairname   string       `json:"pairname"`
	Entrytime  string       `json:"entrytime"`
	Position   string       `json:"position"`
	Leverage   int8         `json:"leverage"`
	Outtime    string       `json:"outtime"`
	Entryprice float64      `json:"entryprice"`
	Quantity   float64      `json:"quantity"`
	Endprice   float64      `json:"endprice"`
	Pnl        float64      `json:"pnl"`
	Roe        float64      `json:"roe"`
	SettledAt  sql.NullTime `json:"settled_at"`
	CreatedAt  time.Time    `json:"created_at"`
}

type RankingBoard struct {
	UserID       string    `json:"user_id"`
	ScoreID      string    `json:"score_id"`
	FinalBalance float64   `json:"final_balance"`
	Comment      string    `json:"comment"`
	CreatedAt    time.Time `json:"created_at"`
}

type RecommendHistory struct {
	Recommender string    `json:"recommender"`
	NewMember   string    `json:"new_member"`
	CreatedAt   time.Time `json:"created_at"`
}

type Session struct {
	SessionID    string    `json:"session_id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type UsedToken struct {
	ScoreID         string    `json:"score_id"`
	UserID          string    `json:"user_id"`
	MetamaskAddress string    `json:"metamask_address"`
	CreatedAt       time.Time `json:"created_at"`
}

type User struct {
	ID                int64          `json:"id"`
	UserID            string         `json:"user_id"`
	OauthUid          sql.NullString `json:"oauth_uid"`
	Nickname          string         `json:"nickname"`
	HashedPassword    sql.NullString `json:"hashed_password"`
	Email             string         `json:"email"`
	MetamaskAddress   sql.NullString `json:"metamask_address"`
	PhotoUrl          sql.NullString `json:"photo_url"`
	PracBalance       float64        `json:"prac_balance"`
	CompBalance       float64        `json:"comp_balance"`
	WmoiBalance       float64        `json:"wmoi_balance"`
	RecommenderCode   string         `json:"recommender_code"`
	CreatedAt         time.Time      `json:"created_at"`
	LastAccessedAt    sql.NullTime   `json:"last_accessed_at"`
	PasswordChangedAt time.Time      `json:"password_changed_at"`
	AddressChangedAt  sql.NullTime   `json:"address_changed_at"`
}

type VerifyEmail struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	SecretCode string    `json:"secret_code"`
	IsUsed     bool      `json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type WmoiMintingHistory struct {
	ToUser    string    `json:"to_user"`
	Amount    int64     `json:"amount"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"method"`
	Giver     string    `json:"giver"`
}
