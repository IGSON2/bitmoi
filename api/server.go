package api

import (
	"bitmoi/contract"
	db "bitmoi/db/sqlc"
	"bitmoi/token"
	"bitmoi/utilities"
	bitmoicommon "bitmoi/utilities/common"
	"bitmoi/worker"
	"log"

	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

const (
	finalstage  = 10
	competition = "competition"
	practice    = "practice"
)

var (
	errNotAuthenticated = errors.New("this service requires authentication in competition mode")
)

type Server struct {
	config          *utilities.Config       // 환경 구성 요소
	googleOauthCfg  *oauth2.Config          // OAuth2.0 인증 구성 요소
	logger          *zerolog.Logger         // 로그 출력
	store           db.Store                // DB 커넥션
	router          *fiber.App              // 각 Endpoint별 router 집합
	tokenMaker      *token.PasetoMaker      // 인증, 인가에 필요한 토큰 생성 및 검증
	pairs           []string                // 제공하는 차트 종목 목록
	taskDistributor worker.TaskDistributor  // Redis 테스크 큐
	erc20Contract   *contract.ERC20Contract // MOI 스마트 컨트랙트
	nextUnlockDate  time.Time               // 경매 종료 일자
	s3Uploader      *s3.S3                  // S3 FullAccess role이 부여된 사용자
	exitCh          chan struct{}           // 서버 종료 시그널을 수신할 채널
	faucetTimeouts  map[string]int64        // 무료 토큰 수령자 별 재요청 제한시간 맵
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(c *utilities.Config, s db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	log.Printf("Environment: %s", c.Environment)
	tm, err := token.NewPasetoTokenMaker(c.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker : %w", err)
	}

	erc20, err := contract.InitErc20Contract(c.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("cannot init erc20 contract : %w", err)
	}

	s3Uploader, err := NewS3Uploader(c)
	if err != nil {
		return nil, fmt.Errorf("cannot init s3 uploader : %w", err)
	}

	server := &Server{
		config:          c,
		googleOauthCfg:  NewGoogleOauthConfig(c),
		store:           s,
		tokenMaker:      tm,
		taskDistributor: taskDistributor,
		erc20Contract:   erc20,
		s3Uploader:      s3Uploader,
		exitCh:          make(chan struct{}),
		nextUnlockDate:  time.Now().Add(c.BiddingDuration),
		faucetTimeouts:  make(map[string]int64),
	}

	ps, err := server.store.GetAllPairsInDB1H(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	router := fiber.New(fiber.Config{})

	lgr := server.createLoggerMiddleware()

	if c.Environment == bitmoicommon.EnvProduction {
		router.Use(createNewOriginMiddleware(), lgr)
	} else {
		router.Use(
			logger.New(logger.Config{Format: "[${ip}]:${port} ${time} ${status} - ${method} ${path} - ${latency}\n"}),
			cors.New(cors.Config{
				AllowOrigins: "*",
			}),
		)
	}

	nomalGroup := router.Group("/basic", createNewLimitMiddleware(30, server.logger))
	nomalGroup.Get("/practice", server.getPracticeChart)
	nomalGroup.Post("/practice", server.postPracticeScore)
	nomalGroup.Get("/interval", server.getAnotherInterval)
	nomalGroup.Get("/score/:nickname", server.getUserScoreSummary)
	nomalGroup.Post("/user/login", server.loginUser)
	nomalGroup.Get("/user/checkId", server.checkID)
	nomalGroup.Get("/user/checkNickname", server.checkNickname)
	nomalGroup.Post("/reissueAccess", server.reissueAccessToken)
	nomalGroup.Post("/verifyToken", server.verifyToken)
	nomalGroup.Get("/login/google", server.GoogleLogin)
	nomalGroup.Get("/login/kakao", server.KakaoLogin)
	nomalGroup.Get("/oauth/:req_url", server.GetLoginURL)
	nomalGroup.Get("/rank", server.getRank)

	authGroup := router.Group("/auth", createNewLimitMiddleware(30, server.logger), authMiddleware(server.tokenMaker))
	authGroup.Get("/myscore", server.myscore)
	authGroup.Get("/checkAttendance", server.checkAttendance)
	authGroup.Post("/user/recommender", server.rewardRecommender)
	authGroup.Get("/user/accumulation", server.getAccumulationHist)
	authGroup.Put("/user/nickname", server.updateNickname)
	authGroup.Put("/user/profile", server.updateProfileImg)
	authGroup.Get("/user/wmoi-transactions", server.getWmoiMintingHist)

	upperLimitedGroup := router.Group("/intermediate", authMiddleware(server.tokenMaker), createNewLimitMiddleware(100, server.logger))
	upperLimitedGroup.Post("/", server.getImdChart)
	upperLimitedGroup.Post("/init", server.initImdScore)
	upperLimitedGroup.Post("/close", server.closeImdScore)
	upperLimitedGroup.Get("/interval", server.getImdInterval)
	upperLimitedGroup.Put("/settle", server.SettleImdScore)

	adminGroup := router.Group("/admin", adminAuthMiddleware(c.AdminID, server.tokenMaker), createNewLimitMiddleware(100, server.logger))
	adminGroup.Get("/users", server.GetUsers)
	adminGroup.Get("/scores", server.GetScoresInfo)
	adminGroup.Get("/usdp", server.GetPracUsdpInfo)
	adminGroup.Post("/usdp", server.SetPracUsdpInfo)
	adminGroup.Get("/token", server.GetTokenInfo)
	adminGroup.Get("/referral", server.GetReferralInfo)

	advancedGroup := router.Group("/advanced", createNewLimitMiddleware(500, server.logger))
	advancedGroup.Get("/practice", server.GetAdvancedPractice)

	// websocketGroup := router.Group("/ws", createWebsocketMiddleware())
	// websocketGroup.Get("/", websocket.New(server.websocketTest))

	// Only bitmoi PC
	nomalGroup.Get("/moreinfo", server.moreinfo)
	nomalGroup.Post("/user", server.createUser)
	nomalGroup.Get("/user/verifyEmail", server.verifyEmail)
	nomalGroup.Get("/nextBidUnlock", server.getNextUnlockDate)
	nomalGroup.Get("/highestBidder", server.getHighestBidder)
	nomalGroup.Get("/selectedBidder", server.getSelectedBidder)

	authGroup.Get("/competition", server.getCompetitionChart)
	authGroup.Post("/competition", server.postCompetitionScore)
	authGroup.Post("/freeToken", server.sendFreeErc20)
	authGroup.Post("/user/address", server.updateMetamaskAddress)
	authGroup.Post("/bidToken", server.bidToken)

	server.router = router

	go server.BiddingLoop()

	return server, nil
}

// Listen enables the server to listen for incoming requests.
func (s *Server) Listen() error {
	errCh := make(chan error)
	go func(ch chan<- error) {
		errCh <- s.router.Listen(s.config.HTTPAddress)
	}(errCh)

	select {
	case err := <-errCh:
		return err
	case <-s.exitCh:
		return ErrClosedBiddingLoop
	}
}

// getPracticeChart godoc
// @Summary      연습모드에서 제공할 차트를 불러옵니다.
// @Tags         chart
// @Param names query string false "제외할 USDT페어들을 쉼표로 구분하여 전달합니다."
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /practice [get]
func (s *Server) getPracticeChart(c *fiber.Ctx) error {
	var oc *OnePairChart
	r := new(CandlesRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}

	history := utilities.SplitPairnames(r.Names)
	if len(history) >= finalstage {
		return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
	}
	for i := 1; ; i++ {
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err = s.makeChartToRef(db.OneH, nextPair, practice, c.Context())
		if err != nil || oc == nil {
			if err == ErrShortRange {
				s.logger.Warn().Str("parname", nextPair).Msgf("%s, Cnt: %d", err.Error(), i)
				if i > 10 {
					s.logger.Error().Str("parname", nextPair).Msg("DB is not ready. not another longer chart.")
					return c.Status(fiber.StatusInternalServerError).SendString("cannot make chart.")
				}
				continue
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		break
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}

// postPracticeScore godoc
// @Summary		연습모드에서 작성한 주문을 제출합니다.
// @Tags		score
// @Accept		json
// @Produce		json
// @Param		order	body		api.ScoreRequest	true	"주문 정보"
// @Success		200		{object}	api.ScoreResponse
// @Router       /basic/practice [post]
func (s *Server) postPracticeScore(c *fiber.Ctx) error {
	var PracticeOrder ScoreRequest
	err := c.BodyParser(&PracticeOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	// errs := utilities.ValidateStruct(PracticeOrder)
	// if errs != nil {
	// 	return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	// }

	err = validateOrderRequest(&PracticeOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	r, err := s.createPracResult(&PracticeOrder, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(r)
}

// deprecated temporarily
func (s *Server) getCompetitionChart(c *fiber.Ctx) error {
	var oc *OnePairChart

	r := new(CandlesRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	history := utilities.SplitPairnames(r.Names)

	if len(history) >= finalstage {
		return c.Status(fiber.StatusBadRequest).SendString("invalid current stage")
	}
	for i := 1; ; i++ {
		nextPair := utilities.FindDiffPair(s.pairs, history)
		oc, err = s.makeChartToRef(db.OneH, nextPair, competition, c.Context())
		if err != nil || oc == nil {
			if err == ErrShortRange {
				s.logger.Warn().Str("parname", nextPair).Msgf("%s, Cnt: %d", err.Error(), i)
				if i > 10 {
					s.logger.Error().Str("parname", nextPair).Msg("DB is not ready. not another longer chart.")
					return c.Status(fiber.StatusInternalServerError).SendString("cannot make chart.")
				}
				continue
			}
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		break
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}

// deprecated temporarily
func (s *Server) postCompetitionScore(c *fiber.Ctx) error {
	var CompetitionOrder ScoreRequest
	err := c.BodyParser(&CompetitionOrder)
	if err != nil || CompetitionOrder.Mode != competition {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("%s, mode : %s", err, CompetitionOrder.Mode))
	}

	errs := utilities.ValidateStruct(CompetitionOrder)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}

	err = validateOrderRequest(&CompetitionOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userID := c.Locals(authorizationPayloadKey).(*token.Payload).UserID
	user, err := s.store.GetUser(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot find user")
	}

	var hash *common.Hash
	switch {
	case CompetitionOrder.Stage == 1:
		if CompetitionOrder.Balance != defaultBalance {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("default balance must be %.0f", defaultBalance))
		}
		hash, err = s.spendErc20OnComp(c, CompetitionOrder.ScoreId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	case CompetitionOrder.Stage > 1:
		prevScore, err := s.store.GetCompScoresByStage(c.Context(), db.GetCompScoresByStageParams{
			UserID:  CompetitionOrder.UserId,
			ScoreID: CompetitionOrder.ScoreId,
			Stage:   CompetitionOrder.Stage - 1,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get stage %02d. err: %w", CompetitionOrder.Stage-1, err).Error())
		}
		if prevScore.Stage != CompetitionOrder.Stage-1 {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid stage number")
		}
		if user.CompBalance < (math.Floor(CompetitionOrder.Balance*10) / 10) {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Invalid balance. expected: %.4f, actual: %.4f", user.CompBalance, CompetitionOrder.Balance))
		}
	}

	compResult, err := s.createCompResult(&CompetitionOrder, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = s.insertScore(&CompetitionOrder, compResult.Score, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if hash != nil {
		return c.Status(fiber.StatusOK).JSON(ScoreResponseWithHash{ScoreResponse: *compResult, TxHash: hash})
	}

	return c.Status(fiber.StatusOK).JSON(compResult)
}

// getAnotherInterval godoc
// @Summary      다른 시간단위의 차트를 불러옵니다. 연습, 경쟁 두 모드 모두 지원합니다.
// @Tags         chart
// @Param 		 anotherIntervalRequest query api.AnotherIntervalRequest true "새로운 시간단위에 대한 요청 정보"
// @param 		 Authorization header string false "Authorization"
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /basic/interval [get]
func (s *Server) getAnotherInterval(c *fiber.Ctx) error {
	r := new(AnotherIntervalRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	if r.Mode == competition {
		if err := checkAuthorization(c, s.tokenMaker); err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("%s err: %s", errNotAuthenticated, err))
		}
	}
	oc, err := s.sendAnotherInterval(r, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}

// myscore godoc
// @Summary      사용자의 경쟁모드 주문 채결 내역을 불러옵니다.
// @Tags         score
// @Param		 page query int true "페이지 번호"
// @Param		 mode query string true "모드"
// @param		 Authorization header string true "Authorization"
// @Produce      json
// @Success      200  {array}  db.CompScore
// @Success      200  {array}  db.PracScore
// @Router       /auth/myscore [get]
func (s *Server) myscore(c *fiber.Ctx) error {
	r := new(MyscoreRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)

	// 추상화 필요
	if r.Mode == practice {
		scores, err := s.getMyPracScores(payload.UserID, int32(r.Page), c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(scores)
	}

	scores, err := s.getMyCompScores(payload.UserID, int32(r.Page), c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}

// deprecated
func (s *Server) moreinfo(c *fiber.Ctx) error {
	r := new(MoreInfoRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	scores, err := s.getCompScoresByScoreID(r.ScoreId, r.UserId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}
