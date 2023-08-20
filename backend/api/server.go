package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	bitmoicommon "bitmoi/backend/utilities/common"
	"bitmoi/backend/worker"

	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
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
	config          *utilities.Config
	store           db.Store
	router          *fiber.App
	tokenMaker      *token.PasetoMaker
	pairs           []string
	taskDistributor worker.TaskDistributor
	erc20Contract   *contract.ERC20Contract
	nextUnlockDate  time.Time
	s3Uploader      *s3.S3
	exitCh          chan struct{}
}

func NewServer(c *utilities.Config, s db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
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
		store:           s,
		tokenMaker:      tm,
		taskDistributor: taskDistributor,
		erc20Contract:   erc20,
		s3Uploader:      s3Uploader,
		exitCh:          make(chan struct{}),
	}

	ps, err := server.store.GetAllParisInDB(context.Background())
	if err != nil {
		return nil, err
	}
	server.pairs = ps

	router := fiber.New(fiber.Config{})

	router.Use(allowOriginMiddleware, limiterMiddleware)
	if c.Environment == bitmoicommon.EnvProduction {
		router.Use(server.createLoggerMiddleware())
	}

	router.Get("/practice", server.getPracticeChart)
	router.Post("/practice", server.postPracticeScore)
	router.Get("/interval", server.getAnotherInterval)
	router.Get("/moreinfo", server.moreinfo)
	router.Get("/rank/:page", server.getRank)
	router.Post("/user", server.createUser)
	router.Post("/user/login", server.loginUser)
	router.Get("/user/checkId", server.checkID)
	router.Get("/user/checkNickname", server.checkNickname)
	router.Get("/user/verifyEmail", server.verifyEmail)
	router.Post("/reissueAccess", server.reissueAccessToken)
	router.Post("/verifyToken", server.verifyToken)
	router.Get("/nextBidUnlock", server.getNextUnlockDate)
	router.Get("/highestBidder", server.getHighestBidder)

	authGroup := router.Group("/", authMiddleware(server.tokenMaker))
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

func (s *Server) Listen() error {
	defer func() {
		s.exitCh <- struct{}{}
	}()
	return s.router.Listen(s.config.HTTPAddress)
}

// getPracticeChart godoc
// @Summary      Get practice chart
// @Description  Get non encrypted chart for practice
// @Tags         chart
// @Param names query string false "Comma-separated list of pair names"
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /practice [get]
func (s *Server) getPracticeChart(c *fiber.Ctx) error {
	r := new(CandlesRequest)
	err := c.QueryParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
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
//
//		@Summary		Post score
//		@Description	Post score of practice mode
//		@Tags			score
//		@Accept			json
//		@Produce		json
//		@Param			order	body		api.ScoreRequest	true	"order inforamtion"
//		@Success		200		{object}	api.ScoreResponse
//	 @Router       /practice [post]
func (s *Server) postPracticeScore(c *fiber.Ctx) error {
	var PracticeOrder ScoreRequest
	err := c.BodyParser(&PracticeOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	errs := utilities.ValidateStruct(PracticeOrder)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}

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
// @Summary      Get competition chart
// @Description  Get encrypted chart for competition
// @Tags         chart
// @Param names query string false "Comma-separated list of pair names"
// @param Authorization header string true "Authorization"
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /competition [get]
func (s *Server) getCompetitionChart(c *fiber.Ctx) error {
	r := new(CandlesRequest)
	err := c.QueryParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
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
//
//	@Summary		Post score
//	@Description	Post score of competition mode
//	@Tags			score
//	@Accept			json
//	@Produce		json
//	@Param			order	body		api.ScoreRequest	true	"order inforamtion"
//	@param Authorization header string true "Authorization"
//	@Success		200		{object}	api.ScoreResponse
//	@Router       /competition [post]
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
// @Summary      Get another interval chart
// @Description  Get another interval chart for both of practice and competition
// @Tags         chart
// @Param anotherIntervalRequest query api.AnotherIntervalRequest true "query has interval,identifier,mode,stage"
// @param Authorization header string false "Authorization"
// @Produce      json
// @Success      200  {object}  api.OnePairChart
// @Router       /interval [get]
func (s *Server) getAnotherInterval(c *fiber.Ctx) error {
	r := new(AnotherIntervalRequest)
	err := c.QueryParser(r)
	if errs := utilities.ValidateStruct(*r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
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
// @Summary      Get myscore
// @Description  Get history of my score
// @Tags         score
// @Param page path int true "page number"
// @param Authorization header string true "Authorization"
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
// @Summary      Get rank
// @Description  Get ranking chart
// @Tags         rank
// @Param page path int true "Page number"
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
// @Summary      Post rank
// @Description  Post rank to ranking chart
// @Tags         rank
// @Param rankInsertRequest body api.RankInsertRequest true "socre id and comment"
// @param Authorization header string true "Authorization"
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
// @Summary      Get moreInfo
// @Description  Get more infomation of rank
// @Tags         rank
// @Param moreInfoRequest query api.MoreInfoRequest true "query contains user id and score id"
// @Produce      json
// @Success      200  {array}  db.Score
// @Router       /moreinfo [get]
func (s *Server) moreinfo(c *fiber.Ctx) error {
	r := new(MoreInfoRequest)
	if err := c.QueryParser(r); err != nil {
		return err
	}
	errs := utilities.ValidateStruct(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
	}
	scores, err := s.getScoresByScoreID(r.ScoreId, r.UserId, c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(scores)
}
