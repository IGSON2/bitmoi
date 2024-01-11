package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

const (
	long                    = "LONG"
	short                   = "SHORT"
	defaultBalance  float64 = 1000
	rankRows                = 100
	myscoreRows             = 15
	rewardReceivers         = 3
)

var (
	ErrNotUpdatedScore    = errors.New("failed to update rank due to low score")
	ErrLiquidation        = errors.New("you were liquidated due to insufficient balance")
	ErrInvalidStageLength = errors.New("insufficient number of stages cleared")
	rewardRates           = []float64{0.5, 0.3, 0.2}
)

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
			ScoreID:       req.GetScoreID(),
			UserID:        req.GetUserID(),
			Stage:         req.GetStage(),
			Pairname:      res.GetPairName(),
			Entrytime:     res.GetEntryTime(),
			Position:      position,
			Leverage:      req.GetLeverage(),
			Outtime:       res.GetOutTime(),
			Entryprice:    req.GetEntryPrice(),
			Quantity:      req.GetQuantity(),
			Endprice:      res.GetEndPrice(),
			Pnl:           res.GetPnl(),
			Roe:           res.GetRoe(),
			RemainBalance: common.FloorDecimal(req.GetBalance() + res.GetPnl()),
		})
	case competition:
		_, err = s.store.InsertCompScore(ctx, db.InsertCompScoreParams{
			ScoreID:       req.GetScoreID(),
			UserID:        req.GetUserID(),
			Stage:         req.GetStage(),
			Pairname:      res.GetPairName(),
			Entrytime:     res.GetEntryTime(),
			Position:      position,
			Leverage:      req.GetLeverage(),
			Outtime:       res.GetOutTime(),
			Entryprice:    req.GetEntryPrice(),
			Quantity:      req.GetQuantity(),
			Endprice:      res.GetEndPrice(),
			Pnl:           res.GetPnl(),
			Roe:           res.GetRoe(),
			RemainBalance: common.FloorDecimal(req.GetBalance() + res.GetPnl()),
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
		_, err = s.store.UpdatePracScore(ctx, db.UpdatePracScoreParams{
			Outtime:       res.GetOutTime(),
			Endprice:      res.GetEndPrice(),
			Pnl:           res.GetPnl(),
			Roe:           res.GetRoe(),
			RemainBalance: common.FloorDecimal(req.GetBalance() + res.GetPnl()),
			UserID:        req.GetUserID(),
			ScoreID:       req.GetScoreID(),
			Stage:         req.GetStage(),
		})
	case competition:
		_, err = s.store.UpdateCompcScore(ctx, db.UpdateCompcScoreParams{
			Pairname:      res.GetPairName(),
			Entrytime:     res.GetEntryTime(),
			Outtime:       res.GetOutTime(),
			Entryprice:    req.GetEntryPrice(),
			Endprice:      res.GetEndPrice(),
			Pnl:           res.GetPnl(),
			Roe:           res.GetRoe(),
			RemainBalance: common.FloorDecimal(req.GetBalance() + res.GetPnl()),
			UserID:        req.GetUserID(),
			ScoreID:       req.GetScoreID(),
			Stage:         req.GetStage(),
		})
	default:
		err = fmt.Errorf("invalid mode: %s", req.GetMode())
	}

	return err
}

func (s *Server) getPracScoreToStage(o *ScoreRequest, c *fiber.Ctx) error {
	i, err := s.store.GetPracScoreToStage(c.Context(), db.GetPracScoreToStageParams{
		ScoreID: o.ScoreId,
		UserID:  o.UserId,
		Stage:   o.Stage,
	})
	if err != nil {
		return fmt.Errorf("cannot get score to current stage, err: %w", err)
	}
	totalScore, ok := i.(float64)
	if !ok {
		return fmt.Errorf("cannot assertion totalScore to float, err: %w", err)
	}

	if defaultBalance+totalScore <= 0 {
		return ErrLiquidation
	}
	return nil
}

func (s *Server) getMyPracScores(userId string, pages int32, c *fiber.Ctx) ([]db.PracScore, error) {
	return s.store.GetPracScoresByUserID(c.Context(), db.GetPracScoresByUserIDParams{
		UserID: userId,
		Limit:  myscoreRows,
		Offset: (pages - 1) * myscoreRows,
	})
}

func (s *Server) getPracScoresByScoreID(scoreId, userId string, c context.Context) ([]db.PracScore, error) {
	return s.store.GetPracScoresByScoreID(c, db.GetPracScoresByScoreIDParams{
		ScoreID: scoreId,
		UserID:  userId,
	})
}

func (s *Server) getAllRanks(pages int32, c *fiber.Ctx) ([]db.RankingBoard, error) {
	return s.store.GetAllRanks(context.Background(), db.GetAllRanksParams{
		Limit:  rankRows,
		Offset: (pages - 1) * rankRows,
	})
}

func (s *Server) getRankByUserID(userId string) (db.RankingBoard, error) {
	return s.store.GetRankByUserID(context.Background(), userId)
}

func (s *Server) insertPracScoreToRankBoard(req *RankInsertRequest, user *db.User, c *fiber.Ctx) error {
	length, err := s.store.GetPracStageLenByScoreID(c.Context(), db.GetPracStageLenByScoreIDParams{
		ScoreID: req.ScoreId,
		UserID:  user.UserID,
	})
	if err != nil {
		return err
	} else if length != finalstage {
		return ErrInvalidStageLength
	}

	t, err := s.store.GetPracScoreToStage(c.Context(), db.GetPracScoreToStageParams{
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
				Nickname:     user.Nickname.String,
				PhotoUrl:     user.PhotoUrl.String,
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
		Nickname:     user.Nickname.String,
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
			log.Error().Err(err).Msgf("cannot find user from db. user: %s", t.UserID)
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
