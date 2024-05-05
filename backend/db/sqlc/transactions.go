package db

import (
	"bitmoi/backend/contract"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User User
}

func (s *SqlStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		result.User, err = q.GetUser(ctx, arg.UserID)
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})

	return result, err
}

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SqlStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		_, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.VerifyEmail, err = q.GetVerifyEmails(ctx, GetVerifyEmailsParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		if time.Now().After(result.VerifyEmail.ExpiredAt) {
			return fmt.Errorf("verification has expired")
		}

		// _, err = q.UpdateUserEmailVerified(ctx, UpdateUserEmailVerifiedParams{
		// 	IsEmailVerified: true,
		// 	UserID:          result.VerifyEmail.UserID,
		// })
		// if err != nil {
		// 	return err
		// }

		result.User, err = q.GetUser(ctx, result.VerifyEmail.UserID)
		return err
	})

	return result, err
}

type SpendTokenTxParams struct {
	CreateUsedTokenParams
	Contract *contract.ERC20Contract
	FromAddr common.Address
	Amount   *big.Int
}

type SpendTokenTxResult struct {
	TxHash *common.Hash
}

func (store *SqlStore) SpendTokenTx(ctx context.Context, arg SpendTokenTxParams) (SpendTokenTxResult, error) {
	var result SpendTokenTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		_, err := q.CreateUsedToken(ctx, arg.CreateUsedTokenParams)
		if err != nil {
			return err
		}
		hash, err := arg.Contract.SpendTokens(arg.FromAddr, arg.Amount, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
		result.TxHash = hash
		return err
	})

	return result, err
}

type AppendPracBalanceTxParams struct {
	UserID string
	Amount float64
	Title  string
	Giver  string
	Method string
}

func (store *SqlStore) AppendPracBalanceTx(ctx context.Context, arg AppendPracBalanceTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		_, err := q.CreateAccumulationHist(ctx, CreateAccumulationHistParams{
			ToUser: arg.UserID,
			Amount: arg.Amount,
			Title:  arg.Title,
			Giver:  arg.Giver,
			Method: arg.Method,
		})
		if err != nil {
			return fmt.Errorf("failed to update reward due to cannot create accumulation history. err: %w", err)
		}

		_, err = q.AppendUserPracBalance(ctx, AppendUserPracBalanceParams{
			PracBalance: arg.Amount,
			UserID:      arg.UserID,
		})
		if err != nil {
			return fmt.Errorf("failed to update reward due to cannot update prac balance. err: %w", err)
		}
		return nil
	})
	return err
}

type CheckAttendTxParams struct {
	AppendPracBalanceTxParams
	TodayMidnight time.Time
}

func (store *SqlStore) CheckAttendTx(ctx context.Context, arg CheckAttendTxParams) (float64, error) {

	var lastAccessed time.Time
	err := store.execTx(ctx, func(q *Queries) error {
		user, err := q.GetUser(ctx, arg.UserID)
		if user.UserID == "" || err != nil {
			return fmt.Errorf("failed to attendence due to cannot find user. err: %w", err)
		}

		lastAccessed = user.LastAccessedAt.Time

		_, err = q.UpdateUserLastAccessedAt(ctx, UpdateUserLastAccessedAtParams{
			// Mysql DB의 time_zone이 Asia/Seoul이여도, Default로 생성되는 값에만 적용되고 Update로 변환되는 argument에는 UTC가 적용된다.
			LastAccessedAt: sql.NullTime{Time: time.Now().Add(9 * time.Hour), Valid: true},
			UserID:         arg.UserID,
		})

		if err != nil {
			return fmt.Errorf("failed to attendence due to cannot update last accessed at. err: %w", err)
		}
		return nil
	})

	if err != nil {
		return arg.Amount, err
	}

	if lastAccessed.After(arg.TodayMidnight) {
		return arg.Amount, errors.New("already checked attendance today")
	}

	return arg.Amount, store.AppendPracBalanceTx(ctx, arg.AppendPracBalanceTxParams)
}

type SettleImdScoreTxParams struct {
	UserID string
}

func (store *SqlStore) SettleImdPracScoreTx(ctx context.Context, arg SettleImdScoreTxParams) (float64, error) {
	totalPnl := 0.0

	err := store.execTx(ctx, func(q *Queries) error {
		scores, err := q.GetUnsettledPracScores(ctx, arg.UserID)
		if err != nil {
			return fmt.Errorf("failed to settle immediate score due to cannot get unsettled scores. err: %w", err)
		}
		if len(scores) == 0 {
			return nil
		}

		for _, sc := range scores {
			totalPnl += sc.Pnl
			_, err = q.UpdatePracScoreSettledAt(ctx, UpdatePracScoreSettledAtParams{
				SettledAt: sql.NullTime{Time: time.Now(), Valid: true},
				UserID:    arg.UserID,
				ScoreID:   sc.ScoreID,
			})
			if err != nil {
				return fmt.Errorf("failed to settle immediate score due to cannot update settled at. err: %w, user:%s, score_id: %s", err, arg.UserID, sc.ScoreID)
			}
		}

		_, err = q.AppendUserPracBalance(ctx, AppendUserPracBalanceParams{
			PracBalance: totalPnl,
			UserID:      arg.UserID,
		})

		if err != nil {
			return fmt.Errorf("failed to settle immediate score due to cannot update user balance. err: %w", err)
		}
		return nil
	})
	return totalPnl, err
}

type RewardRecommenderTxParams struct {
	NewMember       string
	RecommenderCode string
}

const (
	RecommendationReward = 10
	RecommendationTitle  = "추천 보상"
)

var ErrRecommenderNotFound = fmt.Errorf("recommender not found")

func (store *SqlStore) RewardRecommenderTx(ctx context.Context, arg RewardRecommenderTxParams) (string, error) {
	recmNick := ""
	err := store.execTx(ctx, func(q *Queries) error {
		recommender, err := store.GetUserByRecommenderCode(ctx, arg.RecommenderCode)
		if err != nil {
			return ErrRecommenderNotFound
		}

		_, err = store.CreateRecommendHistory(ctx, CreateRecommendHistoryParams{
			Recommender: recommender.UserID,
			NewMember:   arg.NewMember,
		})
		if err != nil {
			return fmt.Errorf("failed to reward recommender due to cannot create recommend history. err: %w, recommender: %s, new_member: %s", err, recommender.UserID, arg.NewMember)
		}

		_, err = store.CreateWmoiMintingHist(ctx, CreateWmoiMintingHistParams{
			ToUser: recommender.UserID,
			Amount: RecommendationReward,
			Title:  RecommendationTitle,
		})
		if err != nil {
			return fmt.Errorf("failed to reward recommender due to cannot create wmoi minting history. err: %w, recommender: %s", err, recommender.UserID)
		}

		recmNick = recommender.Nickname
		return nil
	})
	return recmNick, err
}
