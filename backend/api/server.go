package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	bitmoicommon "bitmoi/backend/utilities/common"
	"bitmoi/backend/worker"
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
	config          *utilities.Config // 환경 구성 요소
	oauthConfig     *oauth2.Config    // OAuth2.0 인증 구성 요소
	logger          zerolog.Logger
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
		oauthConfig:     NewOauthConfig(c),
		store:           s,
		tokenMaker:      tm,
		taskDistributor: taskDistributor,
		erc20Contract:   erc20,
		s3Uploader:      s3Uploader,
		exitCh:          make(chan struct{}),
		nextUnlockDate:  time.Now().Add(c.BiddingDuration),
		faucetTimeouts:  make(map[string]int64),
	}

	ps, err := server.store.GetAllParisInDB(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	router := fiber.New(fiber.Config{})
	router.Use(createNewLimitMiddleware())

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

	router.Get("/practice", server.getPracticeChart)
	router.Post("/practice", server.postPracticeScore)
	router.Get("/interval", server.getAnotherInterval)
	router.Get("/moreinfo", server.moreinfo)
	router.Get("/rank/:page", server.getRank)
	// router.Post("/user", server.createUser)
	// router.Post("/user/login", server.loginUser)
	router.Get("/user/checkId", server.checkID)
	router.Get("/user/checkNickname", server.checkNickname)
	// router.Get("/user/verifyEmail", server.verifyEmail)
	router.Post("/reissueAccess", server.reissueAccessToken)
	router.Post("/verifyToken", server.verifyToken)
	router.Get("/nextBidUnlock", server.getNextUnlockDate)
	router.Get("/highestBidder", server.getHighestBidder)
	router.Get("/selectedBidder", server.getSelectedBidder)
	router.Get("/login/callback", server.CallBackLogin)
	router.Get("/oauth", server.GetLoginURL)

	authGroup := router.Group("/", authMiddleware(server.tokenMaker))
	authGroup.Use(createNewLimitMiddleware())

	if c.Environment == bitmoicommon.EnvProduction {
		authGroup.Use(createNewOriginMiddleware(), lgr)
	} else {
		authGroup.Use(
			logger.New(logger.Config{Format: "[${ip}]:${port} ${time} ${status} - ${method} ${path} - ${latency}\n"}),
			cors.New(cors.Config{
				AllowOrigins: "*",
			}),
		)
	}

	authGroup.Use(lgr)

	authGroup.Get("/competition", server.getCompetitionChart)
	authGroup.Post("/competition", server.postCompetitionScore)
	authGroup.Post("/rank", server.postRank)
	authGroup.Get("/myscore/:page", server.myscore)
	authGroup.Post("/freeToken", server.sendFreeErc20)
	authGroup.Post("/user/address", server.updateMetamaskAddress)
	authGroup.Post("/user/profile", server.updateProfileImg)
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
	nextPair := utilities.FindDiffPair(s.pairs, history)
	oc, err := s.makeChartToRef(db.OneH, nextPair, practice, len(history), c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
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
// @Router       /practice [post]
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

// getCompetitionChart godoc
// @Summary      경쟁모드에서 제공할 차트를 불러옵니다.
// @Tags         chart
// @Param 		 names query string false "제외할 USDT페어들을 쉼표로 구분하여 전달합니다."
// @param 		 Authorization header string true "Authorization"
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /competition [get]
func (s *Server) getCompetitionChart(c *fiber.Ctx) error {
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
	nextPair := utilities.FindDiffPair(s.pairs, history)
	oc, err := s.makeChartToRef(db.OneH, nextPair, competition, len(history), c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(oc)
}

// postCompetitionScore godoc
// @Summary		경쟁모드에서 작성한 주문을 제출합니다.
// @Tags		score
// @Accept		json
// @Produce		json
// @Param		order	body		api.ScoreRequest	true	"주문 정보"
// @param		Authorization header string true "Authorization"
// @Success		200		{object}	api.ScoreResponse
// @Router      /competition [post]
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
		prevScore, err := s.store.GetScore(c.Context(), db.GetScoreParams{
			ScoreID: CompetitionOrder.ScoreId,
			Stage:   CompetitionOrder.Stage - 1,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get stage %02d. err: %w", CompetitionOrder.Stage-1, err).Error())
		}
		if prevScore.Stage != CompetitionOrder.Stage-1 {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid stage number")
		}
		if prevScore.RemainBalance < (math.Floor(CompetitionOrder.Balance*10) / 10) {
			return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Invalid balance. expected: %.4f, actual: %.4f", prevScore.RemainBalance, CompetitionOrder.Balance))
		}
	}

	compResult, err := s.createCompResult(&CompetitionOrder, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	err = s.insertUserScore(&CompetitionOrder, compResult.Score, c)
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
// @Router       /interval [get]
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
// @Param		 page path int true "페이지 번호"
// @param		 Authorization header string true "Authorization"
// @Produce      json
// @Success      200  {array}  db.Score
// @Router       /myscore/{page} [get]
func (s *Server) myscore(c *fiber.Ctx) error {
	page, err := c.ParamsInt("page")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)

	scores, err := s.getMyscores(payload.UserID, int32(page), c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}

// getRank godoc
// @Summary      랭크에 등재된 사용자들을 불러옵니다.
// @Tags         rank
// @Param 		 page path int true "페이지 번호"
// @Produce      json
// @Success      200  {array}  db.RankingBoard
// @Router       /rank/{page} [get]
func (s *Server) getRank(c *fiber.Ctx) error {
	pageNum, err := c.ParamsInt("page")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid page number: %d", pageNum))
	}
	if pageNum == 0 {
		pageNum = 1
	}
	ranks, err := s.getAllRanks(int32(pageNum), c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(ranks)
}

// postRank godoc
// @Summary      사용자를 랭크에 등재합니다.
// @Tags         rank
// @Param 		 rankInsertRequest body api.RankInsertRequest true "랭크 등재 요청에 대한 정보"
// @param 		 Authorization header string true "Authorization"
// @Produce      json
// @Success      200
// @Router       /rank [post]
func (s *Server) postRank(c *fiber.Ctx) error {
	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
	}

	var r RankInsertRequest
	err = c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errs := utilities.ValidateStruct(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}

	err = s.insertScoreToRankBoard(&r, &user, c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}

// moreinfo godoc
// @Summary      사용자가 랭크에 등재하며 기입한 추가 정보를 불러옵니다.
// @Tags         rank
// @Param moreInfoRequest query api.MoreInfoRequest true "추가 정보 요청에 대한 정보"
// @Produce      json
// @Success      200  {array}  db.Score
// @Router       /moreinfo [get]
func (s *Server) moreinfo(c *fiber.Ctx) error {
	r := new(MoreInfoRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	scores, err := s.getScoresByScoreID(r.ScoreId, r.UserId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}
