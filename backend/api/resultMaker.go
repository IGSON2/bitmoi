package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) selectResultChart(info *utilities.IdentificationData, waitingTerm int, c *fiber.Ctx) (*CandleData, error) {
	var cdd CandleData

	switch info.Interval {
	case db.OneD:
		candles, err := s.store.Get1dResult(c.Context(), db.Get1dResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.OneD, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hResult(c.Context(), db.Get4hResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.FourH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hResult(c.Context(), db.Get1hResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.OneH, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = (&cs).InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mResult(c.Context(), db.Get15mResultParams{info.Name, int64(info.RefTimestamp), int32(db.CalculateTerm(db.FifM, waitingTerm))})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = (&cs).InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		return nil, ErrGetResultChart
	}

	return &cdd, nil
}
