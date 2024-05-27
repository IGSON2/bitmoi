package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
)

func (s *Server) selectInterChart(info *utilities.IdentificationData, interval string, minTime, maxTime int64, ctx context.Context) (*CandleData, error) {
	cdd := new(CandleData)

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandlesRnage(ctx, db.Get1dCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 1d intermediate chart.")
			return nil, err
		}
		cs := Candles1dSlice(candles)
		cdd = cs.InitCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandlesRnage(ctx, db.Get4hCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 4h intermediate chart.")
			return nil, err
		}
		cs := Candles4hSlice(candles)
		cdd = cs.InitCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandlesRnage(ctx, db.Get1hCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 1h intermediate chart.")
			return nil, err
		}
		cs := Candles1hSlice(candles)
		cdd = cs.InitCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandlesRnage(ctx, db.Get15mCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 15m intermediate chart.")
			return nil, err
		}
		cs := Candles15mSlice(candles)
		cdd = cs.InitCandleData()
	}
	if cdd.PData == nil || cdd.VData == nil {
		s.logger.Debug().Str("name", info.Name).Str("interval", interval).Int64("min", minTime).Int64("max", maxTime).Msg("No intermediate chart data.")
		return nil, sql.ErrNoRows
	}

	return cdd, nil
}
