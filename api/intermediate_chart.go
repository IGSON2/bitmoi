package api

import (
	db "bitmoi/db/sqlc"
	"bitmoi/utilities"
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

func (s *Server) getImdChart(c *fiber.Ctx) error {
	req := new(ImdStepRequest)

	code, errStr := s.validateInterScoreReq(req, c)
	if code != fiber.StatusOK {
		return c.Status(code).SendString(errStr)
	}

	res := new(InterScoreResponse)

	if req.MinTimestamp >= req.MaxTimestamp {
		return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("min timestamp is bigger than max timestamp. min : %d, max : %d", req.MinTimestamp, req.MaxTimestamp))
	}

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)
	err := json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	cdd, err := s.selectInterChart(info, req.ReqInterval, req.MinTimestamp, req.MaxTimestamp, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot select intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
	}

	if !info.IsPracticeMode() {
		cdd.encodeChart(info.PriceFactor, info.VolumeFactor, info.TimeFactor)
	}

	anotherIntvMap := make(map[string]*CandleData)
	for _, it := range db.GetAnotherIntervals(req.ReqInterval) {
		anotherCdd, err := s.selectInterChart(info, it, req.MinTimestamp, req.MaxTimestamp, c.Context())
		if err != nil {
			if err != sql.ErrNoRows {
				s.logger.Error().Str("name", info.Name).Str("interval", it).Int64("min", req.MinTimestamp).Int64("max", req.MaxTimestamp).Msg("cannot select intermediate chart of other intervals to reference timestamp.")
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

	var (
		result    *InterMediateResult
		resultCdd *CandleData
	)

	// 작은단위로
	if oc, ok := anotherIntvMap[req.CurInterval]; !ok || oc == nil {
		resultCdd = cdd
	} else {
		// 큰 단위에서 작은단위로 요청하여 현재의 큰 단위의 새로운 캔들이 없는 경우
		resultCdd = anotherIntvMap[req.CurInterval]
	}

	pCutCdd, err := cutResultChart(resultCdd, req.CurTimestamp, true)
	if err != nil {
		s.logger.Warn().Err(err).Str("name", info.Name).Int64("timestamp", req.CurTimestamp).Msg("cannot cut result chart.")
	}
	result = calcImdResult(pCutCdd, &req.ImdScoreRequest, info)

	// 어뷰징 행위 방지를 위해 스코어는 캔들 요청시 매번 업데이트 되어야 함
	err = s.updateScore(req, result, c.Context())
	if err != nil {
		s.logger.Error().Err(err).Str("user id", req.UserId).Str("score id", req.ScoreId).Msg("cannot update score. Not initialized.")
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update score. err : %s", err.Error()))
	}

	if result.OutTime > 0 {
		err = s.updateUserBalance(req, result, c.Context())
		if err != nil {
			s.logger.Error().Str("user id", req.UserId).Str("score id", req.ScoreId).Str("mode", req.Mode).Msg("cannot update user balance.")
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update user balance. err : %s", err.Error()))
		}
		cutResultChart(resultCdd, result.OutTime, false)
		// AfterScore 계산 후 저장 필요
	}
	res = &InterScoreResponse{
		ResultChart:   cdd,
		Score:         result,
		AnotherCharts: anotherIntvMap,
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) initImdScore(c *fiber.Ctx) error {
	req := new(ImdScoreRequest)
	code, errStr := s.validateInterScoreReq(req, c)
	if code != fiber.StatusOK {
		return c.Status(code).SendString(errStr)
	}

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)
	err := json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	tempResult := calcImdResult(&CandleData{}, req, info)

	_, getScoreErr := s.store.GetPracScore(c.Context(), db.GetPracScoreParams{UserID: req.UserId, ScoreID: req.ScoreId, Pairname: req.Name})

	if getScoreErr == nil {
		s.logger.Error().Str("user id", req.UserId).Str("score id", req.ScoreId).Msg("already exist score.")
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

func (s *Server) closeImdScore(c *fiber.Ctx) error {
	req := new(ImdCloseRequest)
	code, errStr := s.validateInterScoreReq(req, c)
	if code != fiber.StatusOK {
		return c.Status(code).SendString(errStr)
	}

	info := new(utilities.IdentificationData)
	infoByte := utilities.DecryptByASE(req.Identifier)
	err := json.Unmarshal(infoByte, info)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot unmarshal chart identifier. err : %s", err.Error()))
	}

	cdd, err := s.selectInterChart(info, req.ReqInterval, req.MinTimestamp, req.MaxTimestamp, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot select intermediate chart to reference timestamp. name : %s, interval : %s, err : %s", info.Name, req.ReqInterval, err.Error()))
	}

	if len(cdd.PData) != 1 {
		err = errors.New("invalid intermediate chart data")
		s.logger.Error().Err(err).Str("name", info.Name).Int64("min", req.MinTimestamp).Int64("max", req.MaxTimestamp).Msg("cannot get intermediate result 15m chart.")
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	score := calcImdResult(cdd, &req.ImdScoreRequest, info)
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
		AfterScore: s.calcAfterImdResult(resultChart, &req.ImdScoreRequest, info),
	}
	res.AfterScore.ClosedTime -= info.RefTimestamp

	_, err = s.store.InsertPracAfterScore(c.Context(), db.InsertPracAfterScoreParams{
		ScoreID:      req.ScoreId,
		UserID:       req.UserId,
		MaxRoe:       res.AfterScore.MaxRoe,
		MinRoe:       res.AfterScore.MinRoe,
		AfterOuttime: res.AfterScore.ClosedTime,
	})

	if err != nil {
		s.logger.Error().Err(err).Str("user id", req.UserId).Str("score id", req.ScoreId).Msg("cannot insert after score.")
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot insert after score. err : %s", err.Error()))

	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) validateInterScoreReq(req ScoreReqInterface, c *fiber.Ctx) (int, string) {
	err := c.BodyParser(req)
	if err != nil {
		return fiber.StatusBadRequest, fmt.Sprintf("parsing err : %s", err.Error())
	}
	if errs := utilities.ValidateStruct(req); errs != nil {
		return fiber.StatusBadRequest, fmt.Sprintf("validation err : %s", errs.Error())
	}

	user, err := s.store.GetUser(c.Context(), req.GetUserID())
	if err != nil {
		return fiber.StatusUnauthorized, fmt.Sprintf("cannot get user. err : %s", err.Error())
	}
	if req.GetMode() == practice {
		req.SetBalance(user.PracBalance)
	} else {
		req.SetBalance(user.CompBalance)
	}

	if err := validateOrderRequest(req); err != nil {
		return fiber.StatusBadRequest, err.Error()
	}
	return fiber.StatusOK, ""
}

func cutResultChart(cdd *CandleData, targetTimestamp int64, isPast bool) (*CandleData, error) {
	if cdd == nil || len(cdd.PData) == 0 {
		return nil, fmt.Errorf("invalid intermediate chart data. %v", cdd)
	}
	var cuttingCnt int
	pdatas := cdd.PData

	if pLen := len(pdatas); pLen > 0 {
		if isPast {
			for i := pLen - 1; i > 0; i-- {
				if pdatas[i].Time <= targetTimestamp {
					cuttingCnt++
				}
			}
			cdd.PData = cdd.PData[:pLen-cuttingCnt]
			cdd.VData = cdd.VData[:pLen-cuttingCnt]
		} else {
			for i := 0; i < pLen; i++ {
				if pdatas[i].Time > targetTimestamp {
					cuttingCnt++
				}
			}
			cdd.PData = cdd.PData[cuttingCnt:]
			cdd.VData = cdd.VData[cuttingCnt:]
		}
	}
	return cdd, nil
}
