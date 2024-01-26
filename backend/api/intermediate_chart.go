package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AfterScore struct {
	ClosedTime int64   `json:"closed_time"`
	MinRoe     float64 `json:"min_roe"`
	MaxRoe     float64 `json:"max_roe"`
}

type CloseInterScoreResponse struct {
	Score      *InterMediateResult `json:"score"`
	AfterScore *AfterScore         `json:"after_score"`
}

type InterScoreResponse struct {
	ResultChart   *CandleData            `json:"result_chart"`
	Score         *InterMediateResult    `json:"score"`
	AnotherCharts map[string]*CandleData `json:"another_charts"`
}

func (s *Server) getInterMediateChart(c *fiber.Ctx) error {
	req := new(InterStepRequest)
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
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot select intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
		}

		result := calculateInterResult(cdd, &req.InterScoreRequest, info)

		if result.OutTime > 0 {
			err = s.updateScore(req, result, c.Context())
			if err != nil {
				s.logger.Error().Str("user id", req.UserId).Str("score id", req.ScoreId).Int32("stage", req.Stage).Msg("cannot update score. Not initialized.")
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update score. err : %s", err.Error()))
			}
			err = s.updateUserBalance(req, result, c.Context())
			if err != nil {
				s.logger.Error().Str("user id", req.UserId).Str("score id", req.ScoreId).Int32("stage", req.Stage).Str("mode", req.Mode).Msg("cannot update user balance.")
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update user balance. err : %s", err.Error()))
			}
		}

		if !info.IsPracticeMode() {
			cdd.encodeChart(info.PriceFactor, info.VolumeFactor, info.TimeFactor)
		}

		anotherIntvMap := make(map[string]*CandleData)
		for _, it := range db.GetAnotherIntervals(req.ReqInterval) {
			anotherCdd, err := s.selectInterChart(info, it, req.MinTimestamp, req.MaxTimestamp, c.Context())
			if err != nil {
				if err != sql.ErrNoRows {
					s.logger.Error().Str("name", info.Name).Str("interval", it).Int64("min", req.MinTimestamp).Int64("max", req.MaxTimestamp).Msg("cannot select intermediate chart to reference timestamp.")
					return c.Status(fiber.StatusInternalServerError).SendString("cannot make intermediate another chart to reference timestamp.")
				}
			}
			if !info.IsPracticeMode() {
				if anotherCdd.PData != nil && anotherCdd.VData != nil {
					anotherCdd.encodeChart(info.PriceFactor, info.VolumeFactor, info.TimeFactor)
				}
			}
			anotherIntvMap[it] = anotherCdd
		}

		res = &InterScoreResponse{
			ResultChart:   cdd,
			Score:         result,
			AnotherCharts: anotherIntvMap,
		}
	case req.MinTimestamp >= req.MaxTimestamp:
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("min timestamp is bigger than max timestamp. min : %d, max : %d", req.MinTimestamp, req.MaxTimestamp))
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) initIntermediateScore(c *fiber.Ctx) error {
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

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)
	err = json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	tempResult := calculateInterResult(&CandleData{}, req, info)

	_, getScoreErr := s.store.GetPracScore(c.Context(), db.GetPracScoreParams{UserID: req.UserId, ScoreID: req.ScoreId, Stage: req.Stage})

	if getScoreErr == nil {
		s.logger.Error().Str("user id", req.UserId).Str("score id", req.ScoreId).Int32("stage", req.Stage).Msg("already exist score.")
		return c.Status(fiber.StatusBadRequest).SendString("already exist score.")
	}

	if getScoreErr == sql.ErrNoRows {
		err = s.insertScore(req, tempResult, c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot insert score. err : %s", err.Error()))
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get score. err : %s", getScoreErr.Error()))
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) closeIntermediateScore(c *fiber.Ctx) error {
	req := new(InterStepRequest)
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

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)
	err = json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	// 종료 시각 기준 15분봉 하나면 스코어 계산 가능
	cdd, err := s.selectInterChart(info, db.FifM, req.MinTimestamp, req.MaxTimestamp, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot select intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
	}

	if len(cdd.PData) != 1 {
		err = errors.New("invalid intermediate chart data")
		s.logger.Error().Err(err).Str("name", info.Name).Int64("min", req.MinTimestamp).Int64("max", req.MaxTimestamp).Msg("cannot get intermediate result 15m chart.")
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	score := calculateInterResult(cdd, &req.InterScoreRequest, info)
	score.OutTime = cdd.PData[0].Time
	score.EndPrice = cdd.PData[0].Close

	err = s.updateScore(req, score, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update score. err : %s", err.Error()))
	}

	err = s.updateUserBalance(req, score, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update user balance. err : %s", err.Error()))
	}

	// 진입시간 기준으로 1D부터 축소 검색 한다
	candles, err := s.store.Get1dResult(c.Context(), db.Get1dResultParams{Name: info.Name, Time: int64(info.RefTimestamp), Limit: 2000})
	if err != nil {
		s.logger.Error().Err(err).Str("name", info.Name).Int64("time", int64(info.RefTimestamp)).Msg("cannot get 1d result chart.")
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot calculate after result chart. err : %s", err.Error()))
	}
	cs := Candles1dSlice(candles)
	resultChart := cs.InitCandleData()

	res := &CloseInterScoreResponse{
		Score:      score,
		AfterScore: s.calculateAfterInterResult(resultChart, &req.InterScoreRequest, info),
	}
	res.AfterScore.ClosedTime -= info.RefTimestamp
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
		s.logger.Debug().Str("name", info.Name).Str("interval", interval).Int64("min", minTime).Int64("max", maxTime).Msg("No intermediate chart data.")
	}

	return cdd, nil
}
