// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
	"database/sql"
)

type Querier interface {
	AppendUserCompBalance(ctx context.Context, arg AppendUserCompBalanceParams) (sql.Result, error)
	AppendUserPracBalance(ctx context.Context, arg AppendUserPracBalanceParams) (sql.Result, error)
	AppendUserWmoiBalance(ctx context.Context, arg AppendUserWmoiBalanceParams) (sql.Result, error)
	CreateAccumulationHist(ctx context.Context, arg CreateAccumulationHistParams) (sql.Result, error)
	CreateBiddingHistory(ctx context.Context, arg CreateBiddingHistoryParams) (sql.Result, error)
	CreateRecommendHistory(ctx context.Context, arg CreateRecommendHistoryParams) (sql.Result, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (sql.Result, error)
	CreateUsedToken(ctx context.Context, arg CreateUsedTokenParams) (sql.Result, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error)
	CreateVerifyEmail(ctx context.Context, arg CreateVerifyEmailParams) (sql.Result, error)
	CreateWmoiMintingHist(ctx context.Context, arg CreateWmoiMintingHistParams) (sql.Result, error)
	DeletePairs15m(ctx context.Context, name string) (sql.Result, error)
	DeletePairs1d(ctx context.Context, name string) (sql.Result, error)
	DeletePairs1h(ctx context.Context, name string) (sql.Result, error)
	DeletePairs4h(ctx context.Context, name string) (sql.Result, error)
	DeletePairs5m(ctx context.Context, name string) (sql.Result, error)
	DeletePracScore(ctx context.Context, arg DeletePracScoreParams) (sql.Result, error)
	DeleteUser(ctx context.Context, userID string) (sql.Result, error)
	Get15mCandles(ctx context.Context, arg Get15mCandlesParams) ([]Candles15m, error)
	Get15mCandlesRnage(ctx context.Context, arg Get15mCandlesRnageParams) ([]Candles15m, error)
	Get15mMinMaxTime(ctx context.Context, name string) (Get15mMinMaxTimeRow, error)
	Get15mResult(ctx context.Context, arg Get15mResultParams) ([]Candles15m, error)
	Get15mVolSumPriceAVG(ctx context.Context, arg Get15mVolSumPriceAVGParams) (Get15mVolSumPriceAVGRow, error)
	Get1dCandles(ctx context.Context, arg Get1dCandlesParams) ([]Candles1d, error)
	Get1dCandlesRnage(ctx context.Context, arg Get1dCandlesRnageParams) ([]Candles1d, error)
	Get1dMinMaxTime(ctx context.Context, name string) (Get1dMinMaxTimeRow, error)
	Get1dResult(ctx context.Context, arg Get1dResultParams) ([]Candles1d, error)
	Get1dVolSumPriceAVG(ctx context.Context, arg Get1dVolSumPriceAVGParams) (Get1dVolSumPriceAVGRow, error)
	Get1hCandles(ctx context.Context, arg Get1hCandlesParams) ([]Candles1h, error)
	Get1hCandlesRnage(ctx context.Context, arg Get1hCandlesRnageParams) ([]Candles1h, error)
	Get1hEntryTimestamp(ctx context.Context, arg Get1hEntryTimestampParams) (int64, error)
	Get1hMinMaxTime(ctx context.Context, name string) (Get1hMinMaxTimeRow, error)
	Get1hResult(ctx context.Context, arg Get1hResultParams) ([]Candles1h, error)
	Get1hVolSumPriceAVG(ctx context.Context, arg Get1hVolSumPriceAVGParams) (Get1hVolSumPriceAVGRow, error)
	Get4hCandles(ctx context.Context, arg Get4hCandlesParams) ([]Candles4h, error)
	Get4hCandlesRnage(ctx context.Context, arg Get4hCandlesRnageParams) ([]Candles4h, error)
	Get4hMinMaxTime(ctx context.Context, name string) (Get4hMinMaxTimeRow, error)
	Get4hResult(ctx context.Context, arg Get4hResultParams) ([]Candles4h, error)
	Get4hVolSumPriceAVG(ctx context.Context, arg Get4hVolSumPriceAVGParams) (Get4hVolSumPriceAVGRow, error)
	Get5mCandles(ctx context.Context, arg Get5mCandlesParams) ([]Candles5m, error)
	Get5mCandlesRnage(ctx context.Context, arg Get5mCandlesRnageParams) ([]Candles5m, error)
	Get5mMinMaxTime(ctx context.Context, name string) (Get5mMinMaxTimeRow, error)
	Get5mResult(ctx context.Context, arg Get5mResultParams) ([]Candles5m, error)
	Get5mVolSumPriceAVG(ctx context.Context, arg Get5mVolSumPriceAVGParams) (Get5mVolSumPriceAVGRow, error)
	GetAccumulationHist(ctx context.Context, arg GetAccumulationHistParams) ([]AccumulationHistory, error)
	GetAdminScores(ctx context.Context, arg GetAdminScoresParams) ([]GetAdminScoresRow, error)
	GetAdminUsdpInfo(ctx context.Context, arg GetAdminUsdpInfoParams) ([]GetAdminUsdpInfoRow, error)
	GetAdminUsers(ctx context.Context, arg GetAdminUsersParams) ([]GetAdminUsersRow, error)
	GetAllPairsInDB1D(ctx context.Context) ([]string, error)
	// --------utils----------------
	GetAllPairsInDB1H(ctx context.Context) ([]string, error)
	GetAllRanks(ctx context.Context, arg GetAllRanksParams) ([]RankingBoard, error)
	GetCompScore(ctx context.Context, arg GetCompScoreParams) (CompScore, error)
	GetCompScoreToStage(ctx context.Context, arg GetCompScoreToStageParams) (interface{}, error)
	GetCompScoresByScoreID(ctx context.Context, arg GetCompScoresByScoreIDParams) ([]CompScore, error)
	GetCompScoresByStage(ctx context.Context, arg GetCompScoresByStageParams) (CompScore, error)
	GetCompScoresByUserID(ctx context.Context, arg GetCompScoresByUserIDParams) ([]CompScore, error)
	GetCompStageLenByScoreID(ctx context.Context, arg GetCompStageLenByScoreIDParams) (int64, error)
	GetHighestBidder(ctx context.Context, arg GetHighestBidderParams) (BiddingHistory, error)
	GetHistoryByLocation(ctx context.Context, arg GetHistoryByLocationParams) ([]BiddingHistory, error)
	GetHistoryByUser(ctx context.Context, userID string) ([]BiddingHistory, error)
	GetLastUserID(ctx context.Context) (int64, error)
	GetPracScore(ctx context.Context, arg GetPracScoreParams) (PracScore, error)
	GetPracScoreToStage(ctx context.Context, arg GetPracScoreToStageParams) (interface{}, error)
	GetPracScoresByStage(ctx context.Context, arg GetPracScoresByStageParams) (PracScore, error)
	GetPracScoresByUserID(ctx context.Context, arg GetPracScoresByUserIDParams) ([]PracScore, error)
	GetPracStageLenByScoreID(ctx context.Context, arg GetPracStageLenByScoreIDParams) (int64, error)
	GetRandomUser(ctx context.Context) (User, error)
	GetRankByUserID(ctx context.Context, userID string) (RankingBoard, error)
	GetSession(ctx context.Context, sessionID string) (Session, error)
	GetTopRankers(ctx context.Context, arg GetTopRankersParams) ([]RankingBoard, error)
	GetUnder1YPairs(ctx context.Context, name string) ([]string, error)
	GetUnsettledCompScores(ctx context.Context, userID string) ([]CompScore, error)
	GetUnsettledPracScores(ctx context.Context, userID string) ([]PracScore, error)
	GetUser(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByMetamaskAddress(ctx context.Context, metamaskAddress sql.NullString) (User, error)
	GetUserByNickName(ctx context.Context, nickname string) (User, error)
	GetUserByRecommenderCode(ctx context.Context, recommenderCode string) (User, error)
	GetUserCompScoreSummary(ctx context.Context, nickname string) (GetUserCompScoreSummaryRow, error)
	GetUserLastAccessedAt(ctx context.Context, userID string) (sql.NullTime, error)
	GetUserPracBalance(ctx context.Context, userID string) (float64, error)
	GetUserPracRankByPNL(ctx context.Context, arg GetUserPracRankByPNLParams) ([]GetUserPracRankByPNLRow, error)
	GetUserPracRankByROE(ctx context.Context, arg GetUserPracRankByROEParams) ([]GetUserPracRankByROERow, error)
	GetUserPracScoreSummary(ctx context.Context, nickname string) (GetUserPracScoreSummaryRow, error)
	GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error)
	GetVerifyEmails(ctx context.Context, arg GetVerifyEmailsParams) (VerifyEmail, error)
	GetWmoiMintingHist(ctx context.Context, arg GetWmoiMintingHistParams) ([]WmoiMintingHistory, error)
	Insert15mCandles(ctx context.Context, arg Insert15mCandlesParams) (sql.Result, error)
	Insert1dCandles(ctx context.Context, arg Insert1dCandlesParams) (sql.Result, error)
	Insert1hCandles(ctx context.Context, arg Insert1hCandlesParams) (sql.Result, error)
	Insert4hCandles(ctx context.Context, arg Insert4hCandlesParams) (sql.Result, error)
	Insert5mCandles(ctx context.Context, arg Insert5mCandlesParams) (sql.Result, error)
	InsertCompScore(ctx context.Context, arg InsertCompScoreParams) (sql.Result, error)
	InsertPracAfterScore(ctx context.Context, arg InsertPracAfterScoreParams) (sql.Result, error)
	InsertPracScore(ctx context.Context, arg InsertPracScoreParams) (sql.Result, error)
	InsertRank(ctx context.Context, arg InsertRankParams) (sql.Result, error)
	UpdateCompScoreSettledAt(ctx context.Context, arg UpdateCompScoreSettledAtParams) (sql.Result, error)
	UpdateCompcScore(ctx context.Context, arg UpdateCompcScoreParams) (sql.Result, error)
	UpdatePracScore(ctx context.Context, arg UpdatePracScoreParams) (sql.Result, error)
	UpdatePracScoreSettledAt(ctx context.Context, arg UpdatePracScoreSettledAtParams) (sql.Result, error)
	UpdateUserLastAccessedAt(ctx context.Context, arg UpdateUserLastAccessedAtParams) (sql.Result, error)
	UpdateUserMetamaskAddress(ctx context.Context, arg UpdateUserMetamaskAddressParams) (sql.Result, error)
	UpdateUserNickname(ctx context.Context, arg UpdateUserNicknameParams) (sql.Result, error)
	UpdateUserPhotoURL(ctx context.Context, arg UpdateUserPhotoURLParams) (sql.Result, error)
	UpdateUserRank(ctx context.Context, arg UpdateUserRankParams) (sql.Result, error)
	UpdateUserWmoiBalanceByRecom(ctx context.Context, arg UpdateUserWmoiBalanceByRecomParams) (sql.Result, error)
	UpdateVerifyEmail(ctx context.Context, arg UpdateVerifyEmailParams) (sql.Result, error)
}

var _ Querier = (*Queries)(nil)
