package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type InterScoreResponse struct {
	ResultChart *CandleData         `json:"result_chart"`
	Score       *InterMediateResult `json:"score"`
}

func (s *Server) getInterMediateChart(c *fiber.Ctx) error {
	req := new(InterScoreRequest)
	err := c.BodyParser(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}

	if err := validateOrderRequest(req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	score, getScoreErr := s.store.GetPracScore(c.Context(), db.GetPracScoreParams{UserID: req.UserId, ScoreID: req.ScoreId, Stage: req.Stage})

	res := new(InterScoreResponse)

	switch {
	case req.MinTimestamp < req.MaxTimestamp:
		info := new(utilities.IdentificationData)
		infoByte := utilities.DecryptByASE(req.Identifier)
		err = json.Unmarshal(infoByte, info)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
		}

		cdd, err := s.selectInterChart(info, req.ReqInterval, req.MinTimestamp, req.MaxTimestamp, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot make intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
		}

		result := calculateInterResult(cdd, req, info)

		if getScoreErr != nil {
			if getScoreErr == sql.ErrNoRows {
				err = s.insertScore(req, result, c.Context())
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot insert score. err : %s", err.Error()))
				}
			} else {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get score. err : %s", getScoreErr.Error()))
			}
		}

		if result.OutTime > 0 {
			err = s.updateScore(req, result, c.Context())
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update score. err : %s", err.Error()))
			}
		}

		if !info.IsPracticeMode() {
			cdd.encodeChart(info.PriceFactor, info.VolumeFactor, info.TimeFactor)
		}

		res = &InterScoreResponse{
			ResultChart: cdd,
			Score:       result,
		}
	case req.MinTimestamp == req.MaxTimestamp:
		result := scoreToInterResult(&score)
		// result 계산 후 업데이트 필요
		res = &InterScoreResponse{
			Score: &result,
		}
	case req.MinTimestamp > req.MaxTimestamp:
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("min timestamp is bigger than max timestamp. min : %d, max : %d", req.MinTimestamp, req.MaxTimestamp))
	}

	return c.Status(fiber.StatusOK).JSON(res)
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
