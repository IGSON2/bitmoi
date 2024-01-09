package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) getInterMediateChart(c *fiber.Ctx) error {
	r := new(InterChartRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(r.Identifier)
	err = json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	cdd, err := s.selectInterChart(info, r.ReqInterval, r.MinTimestamp, r.MaxTimestamp, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot make intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, r.ReqInterval, err.Error()))
	}

	var oc = &OnePairChart{
		OneChart:  cdd,
		EntryTime: utilities.EntryTimeFormatter(cdd.PData[0].Time),
	}

	oc.priceFactor = info.PriceFactor
	oc.timeFactor = info.TimeFactor
	oc.volumeFactor = info.VolumeFactor

	if info.PriceFactor != 0 || info.TimeFactor != 0 || info.VolumeFactor != 0 {
		oc.anonymization()
	} else {
		oc.addIdentifier()
	}

	return c.Status(fiber.StatusOK).JSON(oc)
}

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
		return nil, ErrGetResultChart
	}

	return cdd, nil
}
