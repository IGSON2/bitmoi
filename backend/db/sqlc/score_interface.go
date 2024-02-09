package db

import (
	"database/sql"
	"time"
)

type ScoreInterface interface {
	GetScoreID() string
	GetUserID() string
	GetStage() int32
	GetPairname() string
	GetEntrytime() string
	GetPosition() string
	GetLeverage() int32
	GetOuttime() int64
	GetEntryprice() float64
	GetQuantity() float64
	GetEndprice() float64
	GetPnl() float64
	GetRoe() float64
	GetSettledAt() sql.NullTime
	GetCreatedAt() time.Time
}

func (s *PracScore) GetScoreID() string         { return s.ScoreID }
func (s *PracScore) GetUserID() string          { return s.UserID }
func (s *PracScore) GetStage() int32            { return s.Stage }
func (s *PracScore) GetPairname() string        { return s.Pairname }
func (s *PracScore) GetEntrytime() string       { return s.Entrytime }
func (s *PracScore) GetPosition() string        { return s.Position }
func (s *PracScore) GetLeverage() int32         { return s.Leverage }
func (s *PracScore) GetOuttime() int64          { return s.Outtime }
func (s *PracScore) GetEntryprice() float64     { return s.Entryprice }
func (s *PracScore) GetQuantity() float64       { return s.Quantity }
func (s *PracScore) GetEndprice() float64       { return s.Endprice }
func (s *PracScore) GetPnl() float64            { return s.Pnl }
func (s *PracScore) GetRoe() float64            { return s.Roe }
func (s *PracScore) GetSettledAt() sql.NullTime { return s.SettledAt }
func (s *PracScore) GetCreatedAt() time.Time    { return s.CreatedAt }

func (s *CompScore) GetScoreID() string         { return s.ScoreID }
func (s *CompScore) GetUserID() string          { return s.UserID }
func (s *CompScore) GetStage() int32            { return s.Stage }
func (s *CompScore) GetPairname() string        { return s.Pairname }
func (s *CompScore) GetEntrytime() string       { return s.Entrytime }
func (s *CompScore) GetPosition() string        { return s.Position }
func (s *CompScore) GetLeverage() int32         { return s.Leverage }
func (s *CompScore) GetOuttime() int64          { return s.Outtime }
func (s *CompScore) GetEntryprice() float64     { return s.Entryprice }
func (s *CompScore) GetQuantity() float64       { return s.Quantity }
func (s *CompScore) GetEndprice() float64       { return s.Endprice }
func (s *CompScore) GetPnl() float64            { return s.Pnl }
func (s *CompScore) GetRoe() float64            { return s.Roe }
func (s *CompScore) GetSettledAt() sql.NullTime { return s.SettledAt }
func (s *CompScore) GetCreatedAt() time.Time    { return s.CreatedAt }
