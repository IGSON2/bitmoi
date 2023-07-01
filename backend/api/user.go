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
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func convertUserResponse(user db.User) UserResponse {
	return UserResponse{
		UserID:            user.UserID,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
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
		FullName:       req.FullName,
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
				asynq.ProcessIn(10 * time.Second),
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

	rsp := convertUserResponse(txResult.User)
	return c.Status(fiber.StatusOK).JSON(rsp)
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
