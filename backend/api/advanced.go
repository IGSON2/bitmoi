package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AdvancedPracticeQuery struct {
	From       int64  `json:"from" query:"from"`
	To         int64  `json:"to" query:"to"`
	Identifier string `json:"identifier" query:"identifier"`
}

func (s *Server) GetAdvancedPractice(c *fiber.Ctx) error {
	payload := new(AdvancedPracticeQuery)
	err := c.QueryParser(payload)
	if err != nil {
		s.logger.Err(err).Msg("cannot get advanced practice chart due to parsing query failed")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if payload.Identifier == "" {
		nextPair := utilities.FindDiffPair(s.pairs, []string{})
		aOc, err := s.makeAdvancedChartToRef(db.FifM, nextPair, c.Context())
		if err != nil {
			s.logger.Err(err).Msg("cannot get advanced practice chart")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(fiber.StatusOK).JSON(aOc)
	}

	info, err := utilities.DecodeIdentificationData(payload.Identifier)
	if err != nil {
		s.logger.Err(err).Msg("cannot unmarshal chart identifier.")
		return c.Status(fiber.StatusBadRequest).SendString("cannot unmarshal chart identifier.")
	}

	pvData, err := s.selectAdvancedInterChart(info, db.FifM, payload.From, payload.To, c.Context())
	if err != nil {
		s.logger.Err(err).Msg("cannot select intermediate chart to reference timestamp.")
		return c.Status(fiber.StatusInternalServerError).SendString("cannot select intermediate chart to reference timestamp.")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"pvdata": pvData})
}

type AdvancedOnePairChart struct {
	Name       string            `json:"name"`
	PvDatas    []PriceVolumeData `json:"pvdata"`
	Identifier string            `json:"identifier"`
}

func (s *Server) makeAdvancedChartToRef(interval, name string, ctx context.Context) (*AdvancedOnePairChart, error) {

	min, max, err := s.store.SelectMinMaxTime(interval, name, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot count all rows. name : %s, interval : %s, err : %w", name, interval, err)
	}

	refTimestamp := max - calculateRefTimestamp(max-min)
	if refTimestamp == max {
		return nil, ErrShortRange
	}
	pvDatas, err := s.selectAdvancedChart(name, interval, refTimestamp, ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot make chart to reference timestamp. name : %s, interval : %s, err : %w", name, interval, err)
	}

	identifier := utilities.EncrtpByASE(utilities.IdentificationData{
		Name:         name,
		Interval:     interval,
		RefTimestamp: refTimestamp,
	})
	var aOc = &AdvancedOnePairChart{
		Name:       name,
		PvDatas:    pvDatas,
		Identifier: identifier,
	}
	return aOc, nil
}

const oneTimeAdvancedLoad = 10000

func (s *Server) selectAdvancedChart(name, interval string, refTimestamp int64, ctx context.Context) ([]PriceVolumeData, error) {
	var pvDatas []PriceVolumeData

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandles(ctx, db.Get1dCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeAdvancedLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1dSlice(candles)
		pvDatas = (cs).InitAdvancedCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandles(ctx, db.Get4hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeAdvancedLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles4hSlice(candles)
		pvDatas = (cs).InitAdvancedCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandles(ctx, db.Get1hCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeAdvancedLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles1hSlice(candles)
		pvDatas = (cs).InitAdvancedCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandles(ctx, db.Get15mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeAdvancedLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles15mSlice(candles)
		pvDatas = (cs).InitAdvancedCandleData()
	case db.FiveM:
		candles, err := s.store.Get5mCandles(ctx, db.Get5mCandlesParams{Name: name, Time: refTimestamp, Limit: oneTimeStageLoad})
		if err != nil {
			return nil, err
		}
		cs := Candles5mSlice(candles)
		pvDatas = (cs).InitAdvancedCandleData()
	}
	if len(pvDatas) == 0 {
		return nil, ErrGetStageChart
	}
	return pvDatas, nil
}

func (s *Server) selectAdvancedInterChart(info *utilities.IdentificationData, interval string, minTime, maxTime int64, ctx context.Context) ([]PriceVolumeData, error) {
	var pvDatas []PriceVolumeData

	switch interval {
	case db.OneD:
		candles, err := s.store.Get1dCandlesRnage(ctx, db.Get1dCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 1d intermediate chart.")
			return nil, err
		}
		cs := Candles1dSlice(candles)
		pvDatas = cs.InitAdvancedCandleData()
	case db.FourH:
		candles, err := s.store.Get4hCandlesRnage(ctx, db.Get4hCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 4h intermediate chart.")
			return nil, err
		}
		cs := Candles4hSlice(candles)
		pvDatas = cs.InitAdvancedCandleData()
	case db.OneH:
		candles, err := s.store.Get1hCandlesRnage(ctx, db.Get1hCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 1h intermediate chart.")
			return nil, err
		}
		cs := Candles1hSlice(candles)
		pvDatas = cs.InitAdvancedCandleData()
	case db.FifM:
		candles, err := s.store.Get15mCandlesRnage(ctx, db.Get15mCandlesRnageParams{Name: info.Name, Time: minTime, Time_2: maxTime})
		if err != nil {
			s.logger.Error().Err(err).Str("name", info.Name).Int64("min", minTime).Int64("max", maxTime).Msg("cannot get 15m intermediate chart.")
			return nil, err
		}
		cs := Candles15mSlice(candles)
		pvDatas = cs.InitAdvancedCandleData()
	}
	if len(pvDatas) == 0 {
		s.logger.Debug().Str("name", info.Name).Str("interval", interval).Int64("min", minTime).Int64("max", maxTime).Msg("No intermediate chart data.")
		return nil, errors.New("no chart data")
	}

	return pvDatas, nil
}
