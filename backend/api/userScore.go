package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	long                    = "LONG"
	short                   = "SHORT"
	defaultBalance  float64 = 1000
	myscoreRows             = 15
	rewardReceivers         = 3
)

var (
	ErrNotUpdatedScore    = errors.New("failed to update rank due to low score")
	ErrLiquidation        = errors.New("you were liquidated due to insufficient balance")
	ErrInvalidStageLength = errors.New("insufficient number of stages cleared")
	rewardRates           = []float64{0.5, 0.3, 0.2}
)

func (s *Server) updateUserBalance(req ScoreReqInterface, res OrderResultInterface, ctx context.Context) error {
	user, err := s.store.GetUser(ctx, req.GetUserID())
	if err != nil {
		return err
	}
	switch req.GetMode() {
	case practice:
		_, err = s.store.AppendUserPracBalance(ctx, db.AppendUserPracBalanceParams{
			PracBalance: res.GetPnl() - res.GetCommission(),
			UserID:      req.GetUserID(),
		})
		if err != nil {
			return err
		}
	case competition:
		_, err = s.store.AppendUserCompBalance(ctx, db.AppendUserCompBalanceParams{
			CompBalance: user.CompBalance + res.GetPnl() - res.GetCommission(),
			UserID:      req.GetUserID(),
		})
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Server) insertScore(req ScoreReqInterface, res OrderResultInterface, ctx context.Context) error {
	var err error
	var position string

	if req.GetIsLong() {
		position = long
	} else {
		position = short
	}

	switch req.GetMode() {
	case practice:
		_, err = s.store.InsertPracScore(ctx, db.InsertPracScoreParams{
			ScoreID:    req.GetScoreID(),
			UserID:     req.GetUserID(),
			Stage:      req.GetStage(),
			Pairname:   res.GetPairName(),
			Entrytime:  res.GetEntryTime(),
			Position:   position,
			Leverage:   req.GetLeverage(),
			Outtime:    utilities.EntryTimeFormatter(res.GetOutTime()),
			Entryprice: req.GetEntryPrice(),
			Quantity:   req.GetQuantity(),
			Endprice:   res.GetEndPrice(),
			Pnl:        res.GetPnl(),
			Roe:        res.GetRoe(),
		})
	case competition:
		_, err = s.store.InsertCompScore(ctx, db.InsertCompScoreParams{
			ScoreID:    req.GetScoreID(),
			UserID:     req.GetUserID(),
			Stage:      req.GetStage(),
			Pairname:   res.GetPairName(),
			Entrytime:  res.GetEntryTime(),
			Position:   position,
			Leverage:   req.GetLeverage(),
			Outtime:    utilities.EntryTimeFormatter(res.GetOutTime()),
			Entryprice: req.GetEntryPrice(),
			Quantity:   req.GetQuantity(),
			Endprice:   res.GetEndPrice(),
			Pnl:        res.GetPnl(),
			Roe:        res.GetRoe(),
		})
	default:
		err = fmt.Errorf("invalid mode: %s", req.GetMode())
	}

	return err
}

func (s *Server) updateScore(req ScoreReqInterface, res OrderResultInterface, ctx context.Context) error {
	var err error

	switch req.GetMode() {
	case practice:
		pracScore, getErr := s.store.GetPracScore(ctx, db.GetPracScoreParams{
			UserID:   req.GetUserID(),
			ScoreID:  req.GetScoreID(),
			Pairname: req.GetPairName(),
		})
		if getErr != nil {
			return getErr
		}
		if pracScore.Outtime != "" {
			return errors.New("already closed score")
		}
		_, err = s.store.UpdatePracScore(ctx, db.UpdatePracScoreParams{
			Outtime:  utilities.EntryTimeFormatter(res.GetOutTime()),
			Endprice: res.GetEndPrice(),
			Pnl:      res.GetPnl(),
			Roe:      res.GetRoe(),
			UserID:   req.GetUserID(),
			ScoreID:  req.GetScoreID(),
			Stage:    req.GetStage(),
		})
	case competition:
		compScore, getErr := s.store.GetCompScore(ctx, db.GetCompScoreParams{
			UserID:   req.GetUserID(),
			ScoreID:  req.GetScoreID(),
			Pairname: req.GetPairName(),
		})
		if getErr != nil {
			return getErr
		}
		if compScore.Outtime != "" {
			return errors.New("already closed score")
		}
		_, err = s.store.UpdateCompcScore(ctx, db.UpdateCompcScoreParams{
			Pairname:   res.GetPairName(),
			Entrytime:  res.GetEntryTime(),
			Outtime:    utilities.EntryTimeFormatter(res.GetOutTime()),
			Entryprice: req.GetEntryPrice(),
			Endprice:   res.GetEndPrice(),
			Pnl:        res.GetPnl(),
			Roe:        res.GetRoe(),
			UserID:     req.GetUserID(),
			ScoreID:    req.GetScoreID(),
			Pairname_2: req.GetPairName(), //TODO: pairname_2 비효율적
		})
	default:
		err = fmt.Errorf("invalid mode: %s", req.GetMode())
	}
	return err
}

func (s *Server) getMyPracScores(userId string, page int32, c *fiber.Ctx) ([]db.PracScore, error) {
	return s.store.GetPracScoresByUserID(c.Context(), db.GetPracScoresByUserIDParams{
		UserID: userId,
		Limit:  myscoreRows,
		Offset: (page - 1) * myscoreRows,
	})
}

func (s *Server) getMyCompScores(userId string, pages int32, c *fiber.Ctx) ([]db.CompScore, error) {
	return s.store.GetCompScoresByUserID(c.Context(), db.GetCompScoresByUserIDParams{
		UserID: userId,
		Limit:  myscoreRows,
		Offset: (pages - 1) * myscoreRows,
	})
}

type UserScoreSummary struct {
	db.GetUserPracScoreSummaryRow         // 경쟁모드를 포함하는 추상화 인터페이스 필요
	WeeklyRate                    float64 `json:"weekly_winrate"`
}

// getUserScoreSummary godoc
// @Summary      요청한 닉네임을 가진 사용자의 점수 요약을 반환합니다.
// @Tags         rank
// @Param 		 user path string true "User Nickname"
// @Produce      json
// @Success      200  {array}  UserScoreSummary
// @Router       /score/{user} [get]
func (s *Server) getUserScoreSummary(c *fiber.Ctx) error {
	// mode := c.Query("mode", practice)

	nickname, err := url.QueryUnescape(c.Params("nickname"))
	if err != nil {
		s.logger.Error().Err(err).Msg("cannot unescape nickname")
		return c.Status(fiber.StatusBadRequest).SendString("invalid nickname parameter")
	}

	if nickname == "" {
		return c.Status(fiber.StatusBadRequest).SendString("user id is required")
	}

	result, err := s.store.GetUserPracScoreSummary(c.Context(), nickname)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	var weeklyRate float64
	if result.WeeklyWin > 0 || result.WeeklyLose > 0 {
		floatRate := float64(result.WeeklyWin) / float64(result.WeeklyWin+result.WeeklyLose)
		weeklyRate = math.Floor(10000*floatRate) / 100
	}

	return c.Status(fiber.StatusOK).JSON(UserScoreSummary{result, weeklyRate})
}

func (s *Server) getCompScoresByScoreID(scoreId, userId string, c context.Context) ([]db.CompScore, error) {
	return s.store.GetCompScoresByScoreID(c, db.GetCompScoresByScoreIDParams{
		ScoreID: scoreId,
		UserID:  userId,
	})
}

func (s *Server) getRankByUserID(userId string) (db.RankingBoard, error) {
	return s.store.GetRankByUserID(context.Background(), userId)
}

func (s *Server) insertCompScoreToRankBoard(req *RankInsertRequest, user *db.User, c *fiber.Ctx) error {
	length, err := s.store.GetCompStageLenByScoreID(c.Context(), db.GetCompStageLenByScoreIDParams{
		ScoreID: req.ScoreId,
		UserID:  user.UserID,
	})
	if err != nil {
		return err
	} else if length != finalstage {
		return ErrInvalidStageLength
	}

	t, err := s.store.GetCompScoreToStage(c.Context(), db.GetCompScoreToStageParams{
		ScoreID: req.ScoreId,
		UserID:  user.UserID,
		Stage:   finalstage,
	})
	if err != nil {
		return err
	}
	totalScore, ok := t.(float64)
	if !ok {
		return fmt.Errorf("cannot assign totalscore")
	}

	r, err := s.getRankByUserID(user.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err = s.store.InsertRank(c.Context(), db.InsertRankParams{
				UserID:       user.UserID,
				ScoreID:      req.ScoreId,
				Comment:      req.Comment,
				FinalBalance: math.Floor(100*(totalScore+defaultBalance)) / 100,
			})
		}
		return err
	}

	if r.FinalBalance > totalScore {
		return ErrNotUpdatedScore
	}
	_, err = s.store.UpdateUserRank(c.Context(), db.UpdateUserRankParams{
		UserID:       user.UserID,
		ScoreID:      req.ScoreId,
		Comment:      req.Comment,
		FinalBalance: math.Floor(100*(totalScore+defaultBalance)) / 100,
	})
	return err
}

func (s *Server) SendReward(prevUnlocked time.Time) error {

	top3, err := s.store.GetTopRankers(context.Background(), db.GetTopRankersParams{
		CreatedAt: prevUnlocked.UTC(),
		Limit:     3,
		Offset:    0,
	})

	if len(top3) == 0 {
		log.Warn().Msg("no rankers of top3")
		return nil
	}

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Msg("no rankers of top3")
			return nil
		}
		log.Error().Err(err).Msg("cannot get rankers of top3 from db.")
		return err
	}

	spendCnt, err := s.erc20Contract.GetSpendCounts()
	if err != nil {
		log.Error().Err(err).Msg("cannot get spend counts from contract.")
		return err
	}

	tops := make([]ethcommon.Address, rewardReceivers)
	amounts := make([]*big.Int, rewardReceivers)
	for i, t := range top3 {
		user, err := s.store.GetUser(context.Background(), t.UserID)
		if err != nil {
			log.Error().Err(err).Msgf("cannot find user to send reward from db. user: %s", t.UserID)
			return err
		}

		addr := ethcommon.HexToAddress(user.MetamaskAddress.String)
		tops = append(tops, addr)
		amounts[i] = big.NewInt(int64(float64(spendCnt.Uint64()) / rewardRates[i]))
	}

	hash, err := s.erc20Contract.SendReward(tops, amounts, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
	if err != nil {
		log.Error().Err(err).Msgf("cannot send reward.")
		return err
	}
	_, err = s.erc20Contract.WaitAndReturnTxReceipt(hash)
	if err != nil {
		log.Error().Err(err).Msgf("Cannot get receipt of send reward transaction.")
		return err
	}
	return nil
}

type SpendCount struct {
	Count uint64 `json:"spend_count"`
}

func (s *Server) GetSpendCount(c *fiber.Ctx) error {
	cnt, err := s.erc20Contract.GetSpendCounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(SpendCount{Count: cnt.Uint64()})
}

func (s *Server) SettleImdScore(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot find user")
	}

	userID := payload.UserID
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot find user")
	}

	totalPnl, err := s.store.SettleImdPracScoreTx(c.Context(), db.SettleImdScoreTxParams{UserID: userID})
	if err != nil {
		s.logger.Error().Err(err).Msg("cannot settle immediate score")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(struct {
		TotalPnl float64 `json:"total_pnl"`
	}{totalPnl})
}
