package api

import (
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"bitmoi/backend/worker"
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	MetamaskAddress   string    `json:"metamask_address"`
	PracBalance       float64   `json:"prac_balance"`
	CompBalance       float64   `json:"comp_balance"`
	RecommenderCode   string    `json:"recommender_code"`
	LastAccessedAt    time.Time `json:"last_accessed_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func convertUserResponse(user db.User) UserResponse {
	uR := UserResponse{
		UserID:            user.UserID,
		Nickname:          user.Nickname,
		Email:             user.Email,
		PracBalance:       user.PracBalance,
		CompBalance:       user.CompBalance,
		RecommenderCode:   user.RecommenderCode,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
		LastAccessedAt:    user.LastAccessedAt,
	}
	if user.PhotoUrl.Valid {
		uR.PhotoURL = user.PhotoUrl.String
	}
	if user.MetamaskAddress.Valid {
		uR.MetamaskAddress = user.MetamaskAddress.String
	}
	return uR
}

// checkID godoc
// @Summary      ID 중복 여부를 확인합니다.
// @Tags         user
// @Param user_id query string true "체크할 ID"
// @Produce      json
// @Success      200
// @Router       /user/checkId [get]
func (s *Server) checkID(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	user, _ := s.store.GetUser(c.Context(), userID)
	if user.UserID == userID {
		return c.Status(fiber.StatusBadRequest).SendString("user already exist")
	}
	return c.Status(fiber.StatusOK).SendString(userID)
}

// checkNickname godoc
// @Summary      닉네임 중복 여부를 확인합니다.
// @Tags         user
// @Param nickname query string true "체크할 닉네임"
// @Produce      json
// @Success      200
// @Router       /user/checkNickname [get]
func (s *Server) checkNickname(c *fiber.Ctx) error {
	nickname := c.Query("nickname")
	if nickname == "" {
		return c.Status(fiber.StatusBadRequest).SendString("nickname is empty")
	}
	user, _ := s.store.GetUserByNickName(c.Context(), nickname)
	if user.Nickname == nickname {
		return c.Status(fiber.StatusBadRequest).SendString("full name already exist")
	}
	return c.Status(fiber.StatusOK).SendString(nickname)
}

// createUser godoc
// @Summary			사용자를 생성합니다.
// @Tags			user
// @Accept			json
// @Produce			json
// @Param			CreateUserRequest	body		api.CreateUserRequest	true	"유저 생성 요청에 필요한 정보"
// @Success			200
// @Router       	/user [post]
func (s *Server) createUser(c *fiber.Ctx) error {
	req := &CreateUserRequest{
		UserID:   c.FormValue("user_id"),
		Password: c.FormValue("password"),
		Nickname: c.FormValue("nickname"),
		Email:    c.FormValue("email"),
		OauthUid: c.FormValue("oauth_uid"),
	}
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

	var uploadErr error
	var fileURL string

	f, err := c.FormFile(formFileKey)
	if err != nil {
		s.logger.Debug().Msgf("%s dosen't upload profile image.", req.UserID)
	} else {
		fileURL, uploadErr = s.uploadProfileImageToS3(f, req.UserID)
	}

	arg := db.CreateUserParams{
		UserID:         req.UserID,
		HashedPassword: sql.NullString{Valid: true, String: hashedPassword},
		Nickname:       req.Nickname,
		Email:          req.Email,
	}
	if uploadErr == nil {
		arg.PhotoUrl = sql.NullString{String: fileURL, Valid: true}
	} else {
		s.logger.Error().Err(uploadErr).Msgf("cannot upload image to S3. user: %s", req.UserID)
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
		s.deleteObject(req.UserID)
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

// loginUser godoc
// @Summary		사용자의 ID와 PW를 확인하여 토큰을 발급합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		LoginUserRequest	body		api.LoginUserRequest	true	"ID와 PW"
// @Success		200		{object}	api.LoginUserResponse
// @Router      /user/login [post]
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

	// if !user.IsEmailVerified {
	// 	if time.Now().After(user.CreatedAt.Add(s.config.AccessTokenDuration)) {
	// 		taskPayload := &worker.PayloadSendVerifyEmail{
	// 			UserID: user.UserID,
	// 		}
	// 		opts := []asynq.Option{
	// 			asynq.MaxRetry(10),
	// 			asynq.ProcessIn(10 * time.Second),
	// 			asynq.Queue(worker.QueueCritical),
	// 		}
	// 		err = s.taskDistributor.DistributeTaskSendVerifyEmail(c.Context(), taskPayload, opts...)
	// 		return c.Status(fiber.StatusUnauthorized).SendString(fmt.Sprintf("인증 Email을 한 번 더 보내드렸습니다. 인증을 완료해주세요. : %v", err))
	// 	}
	// 	return c.Status(fiber.StatusUnauthorized).SendString("먼저 Email 인증을 완료해주세요.")
	// }

	if err := utilities.CheckPassword(loginReq.Password, user.HashedPassword.String); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("invalid password")
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

// updateMetamaskAddress godoc
// @Summary		사용자의 Metamask 계좌를 업데이트 합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		LoginUserRequest	body		api.MetamaskAddressRequest	true	"Metamask 주소"
// @param 		Authorization header string true "Authorization"
// @Success		200
// @Router      /user/address [post]
func (s *Server) updateMetamaskAddress(c *fiber.Ctx) error {
	r := new(UpdateMetamaskRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get user by authorization payload")
	}

	timeout := user.AddressChangedAt.Time.Add(24 * time.Hour)
	if time.Now().Before(timeout) {
		return c.Status(fiber.StatusUnprocessableEntity).SendString(fmt.Errorf("%s left until next allowance", common.PrettyDuration(time.Until(timeout))).Error())
	}

	_, err = s.store.UpdateUserMetamaskAddress(c.Context(), db.UpdateUserMetamaskAddressParams{
		MetamaskAddress:  sql.NullString{String: r.Addr, Valid: true},
		UserID:           payload.UserID,
		AddressChangedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("cannot update address err:" + err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// updateProfileImg godoc
// @Summary		사용자의 프로필 이미지를 업데이트 합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		image	formData		file	true	"profile image"
// @param 		Authorization header string true "Authorization"
// @Success		200
// @Router      /user/profile [post]
func (s *Server) updateProfileImg(c *fiber.Ctx) error {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get user by authorization payload")
	}

	f, err := c.FormFile(formFileKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get photo image file from context. err: %w", err).Error())
	}

	url, err := s.uploadProfileImageToS3(f, user.UserID)
	if err != nil {
		s.logger.Error().Err(err).Msgf("cannot upload image to S3. user: %s", user.UserID)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	_, err = s.store.UpdateUserPhotoURL(c.Context(), db.UpdateUserPhotoURLParams{
		PhotoUrl: sql.NullString{Valid: true, String: url},
		UserID:   user.UserID,
	})

	if err != nil {
		s.deleteObject(user.UserID)
		s.logger.Error().Err(err).Msgf("cannot update photo url. user: %s", user.UserID)
		return c.Status(fiber.StatusInternalServerError).SendString("cannot update profile image.")
	}

	return c.Status(fiber.StatusOK).SendString(url)
}

// updateNickname godoc
// @Summary		사용자의 닉네임을 업데이트 합니다.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		LoginUserRequest	body		api.UpdateNicknameRequest	true	"New nickname"
// @param 		Authorization header string true "Authorization"
// @Success		200
// @Router      /user/profile [post]
func (s *Server) updateNickname(c *fiber.Ctx) error {
	r := new(UpdateNicknameRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	checkUser, _ := s.store.GetUserByNickName(c.Context(), r.Nickname)
	if checkUser.Nickname == r.Nickname {
		return c.Status(fiber.StatusBadRequest).SendString("nickname already exist")
	}

	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get user by authorization payload")
	}

	_, err = s.store.UpdateUserNickname(c.Context(), db.UpdateUserNicknameParams{
		Nickname: r.Nickname,
		UserID:   user.UserID,
	})

	if err != nil {
		s.logger.Error().Err(err).Msgf("cannot update nickname. user: %s", user.UserID)
		return c.Status(fiber.StatusInternalServerError).SendString("cannot update nickname")
	}

	return c.SendStatus(fiber.StatusOK)
}

const accumulationRows = 10

func (s *Server) getAccumulationHist(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page < 1 {
		return c.Status(fiber.StatusBadRequest).SendString("page must be greater than 0")
	}

	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).SendString("cannot get authorization payload")
	}

	histories, err := s.store.GetAccumulationHist(c.Context(), db.GetAccumulationHistParams{
		ToUser: payload.UserID,
		Title:  "%",
		Limit:  accumulationRows,
		Offset: (int32(page) - 1) * accumulationRows,
	})

	if err != nil {
		s.logger.Error().Err(err).Msgf("cannot get accumulation history. user_id: %s", payload.UserID)
		return c.Status(fiber.StatusInternalServerError).SendString("cannot get accumulation history")
	}

	return c.Status(fiber.StatusOK).JSON(histories)
}
