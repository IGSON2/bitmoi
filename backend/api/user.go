package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/lib/pq"
)

type UserResponse struct {
	UserID            string    `json:"user_id"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	PhotoURL          string    `json:"photo_url"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func convertUserResponse(user db.User) UserResponse {
	uR := UserResponse{
		UserID:            user.UserID,
		Nickname:          user.Nickname,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
	if user.PhotoUrl.String != "" {
		uR.PhotoURL = user.PhotoUrl.String
	}

	return uR
}

func (s *Server) checkID(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	user, _ := s.store.GetUser(c.Context(), userID)
	if user.UserID == userID {
		return c.Status(fiber.StatusBadRequest).SendString("user already exist")
	}
	return c.Status(fiber.StatusOK).SendString(userID)
}

func (s *Server) checkNickname(c *fiber.Ctx) error {
	nickname := c.Query("nickname")
	user, _ := s.store.GetUserByNickName(c.Context(), nickname)
	if user.Nickname == nickname {
		return c.Status(fiber.StatusBadRequest).SendString("full name already exist")
	}
	return c.Status(fiber.StatusOK).SendString(nickname)
}

func (s *Server) createUser(c *fiber.Ctx) error {
	req := new(CreateUserRequest)
	err := c.BodyParser(req)
	if errs := utilities.ValidateStruct(*req); err != nil || errs != nil {
		if errs != nil {
			return c.Status(fiber.StatusBadRequest).SendString(errs.Error())
		}
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hashedPassword, err := utilities.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	arg := db.CreateUserParams{
		UserID:         req.UserID,
		HashedPassword: hashedPassword,
		Nickname:       req.Nickname,
		Email:          req.Email,
	}
	if req.PhotoUrl != "" {
		arg.PhotoUrl = sql.NullString{String: req.PhotoUrl, Valid: true}
	}
	if req.OauthUid != "" {
		arg.OauthUid = sql.NullString{String: req.OauthUid, Valid: true}
	}

	txArg := db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				UserID: user.UserID,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return s.taskDistributor.DistributeTaskSendVerifyEmail(c.Context(), taskPayload, opts...)
		},
	}

	txResult, err := s.store.CreateUserTx(c.Context(), txArg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return c.Status(fiber.StatusForbidden).SendString(err.Error())
			}
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(txResult.User.Email)
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  UserResponse `json:"user"`
}

func (s *Server) loginUser(c *fiber.Ctx) error {
	loginReq := new(LoginUserRequest)
	err := c.BodyParser(&loginReq)
	if errs := utilities.ValidateStruct(loginReq); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	user, err := s.store.GetUser(c.Context(), loginReq.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if !user.IsEmailVerified {
		if time.Now().After(user.CreatedAt.Add(s.config.AccessTokenDuration)) {
			taskPayload := &worker.PayloadSendVerifyEmail{
				UserID: user.UserID,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			err = s.taskDistributor.DistributeTaskSendVerifyEmail(c.Context(), taskPayload, opts...)
			return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("인증 Email을 한 번 더 보내드렸습니다. 인증을 완료해주세요. : %v", err))
		}
		return c.Status(fiber.StatusUnauthorized).SendString("먼저 Email 인증을 완료해주세요.")
	}

	if err := utilities.CheckPassword(loginReq.Password, user.HashedPassword); err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("password is not correct err : %v", err))
	}

	accessToken, accessPayload, err := s.tokenMaker.CreateToken(
		user.UserID,
		s.config.AccessTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		user.UserID,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = s.store.CreateSession(c.Context(), db.CreateSessionParams{
		SessionID:    refreshPayload.SessionID.String(),
		UserID:       user.UserID,
		RefreshToken: refreshToken,
		UserAgent:    string(c.Request().Header.UserAgent()),
		ClientIp:     c.IP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	rsp := LoginUserResponse{
		SessionID:             refreshPayload.SessionID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  convertUserResponse(user),
	}
	return c.Status(fiber.StatusOK).JSON(rsp)

}

func (s *Server) updateUsingToken(c *fiber.Ctx) error {
	r := new(UpdateMetamaskAddrRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	_, err = s.store.CreateUsedToken(c.Context(), db.CreateUsedTokenParams{
		ScoreID:         r.ScoreId,
		UserID:          r.UserID,
		MetamaskAddress: r.MetamaskAddr,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot update metamask address: user:%s addr:%s", r.UserID, r.MetamaskAddr))
	}
	return c.SendStatus(fiber.StatusOK)
}
