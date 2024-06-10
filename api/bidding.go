package api

import (
	"bitmoi/contract"
	db "bitmoi/db/sqlc"
	"bitmoi/token"
	"bitmoi/utilities"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrClosedBiddingLoop = errors.New("bidding loop has been closed")
)

type NextUnlockResponse struct {
	NextUnlock string `json:"next_unlock"`
}

func (s *Server) BiddingLoop() {
	for {
		biddingTimer := time.NewTimer(s.config.BiddingDuration)
		select {
		case <-biddingTimer.C:
			var initErr error
			s.erc20Contract, initErr = contract.InitErc20Contract(s.config.PrivateKey)
			if initErr != nil {
				log.Error().Err(initErr).Msg("Cannot initialize contract")
				s.exitCh <- struct{}{}
			}
			hash, err := s.erc20Contract.UnLockTokens(contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
			if err != nil {
				log.Error().Err(err).Msg("cannot unlock token.")
				s.exitCh <- struct{}{}
			}
			_, err = s.erc20Contract.WaitAndReturnTxReceipt(hash)
			if err != nil {
				log.Error().Err(err).Msgf("Cannot get receipt of unlock token transaction. stop server. hash:%s", hash.Hex())
				s.exitCh <- struct{}{}
			}

			err = s.SendReward(s.nextUnlockDate.Add(-1 * s.config.BiddingDuration))
			if err != nil {
				log.Error().Err(err).Msgf("send reward err occured during unlock.")
				s.exitCh <- struct{}{}
			}

			s.nextUnlockDate = time.Now().Add(s.config.BiddingDuration)
			log.Info().Msgf("Unlock token successfully. hash: %s, next unlock date: %s", hash.Hex(), s.nextUnlockDate.Format("2006-01-02 15:04:05"))
		case <-s.exitCh:
			log.Warn().Err(ErrClosedBiddingLoop)
			return
		}
	}
}

// getNextUnlockDate godoc
// @Summary      경매 마감 일자를 제공합니다.
// @Tags         erc20
// @Success      200 {object} api.NextUnlockResponse "포멧된 일자"
// @Router       /nextBidUnlock [get]
func (s *Server) getNextUnlockDate(c *fiber.Ctx) error {
	kst := time.FixedZone("KST", 9*60*60)
	kstTime := s.nextUnlockDate.In(kst)
	return c.Status(fiber.StatusOK).JSON(NextUnlockResponse{
		NextUnlock: kstTime.Format("2006-01-02T15:04:05")})
}

type HighestBidderResponse struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}

// getHighestBidder godoc
// @Summary      특정 광고 스팟에 가장 높은 가격을 제시한 입찰자를 반환합니다.
// @Tags         erc20
// @Param location query string true "광고 스팟"
// @Success      200 {object} api.HighestBidderResponse "최상위 입찰자"
// @Router       /highestBidder [get]
func (s *Server) getHighestBidder(c *fiber.Ctx) error {
	r := new(GetBidderByLocRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	addr, amt, err := s.erc20Contract.GetHighestBidder(r.Location)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get highest bidder in contract. err: %s", err.Error()))
	}

	user, err := s.store.GetUserByMetamaskAddress(c.Context(), sql.NullString{String: addr.Hex(), Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("cannot find user by matamask address. err: %s", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get user by matamask address. err: %s", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(HighestBidderResponse{UserID: user.UserID, Amount: amt.Int64()})
}

type BidTokenResponse struct {
	ImageURL string       `json:"image_url"`
	Hash     *common.Hash `json:"hash"`
}

// bidToken godoc
// @Summary		광고 스팟에 MOI 토큰을 입찰합니다.
// @Tags		erc20
// @Accept		json
// @Produce		json
// @Param		request	body		api.BidTokenRequest	true	"입찰 금액과 광고 스팟"
// @param 		Authorization header string true "Authorization"
// @Success		200		{object}	api.ScoreResponse
// @Router      /bidToken [post]
func (s *Server) bidToken(c *fiber.Ctx) error {
	receivedAmt := c.FormValue("amount")
	receivedLoc := c.FormValue("location")
	log.Info().Msgf("Amt: %s, Loc: %s", receivedAmt, receivedLoc)
	amt, err := strconv.Atoi(receivedAmt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("cannot convert string to integer err : %s", err.Error()))
	}

	req := BidTokenRequest{
		Amount:   amt,
		Location: c.FormValue("location"),
	}
	if errs := utilities.ValidateStruct(req); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf(" validation err : %s", errs.Error()))
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
	}

	f, err := c.FormFile(formFileKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Ad image dosen't exist. err: %s", err.Error()))
	}

	fileURL, err := s.uploadADImageToS3(f, user.UserID, req.Location)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Cannot upload ad image to S3 bucket. err: %s", err.Error()))
	}

	addr := common.HexToAddress(user.MetamaskAddress.String)

	balance, err := s.erc20Contract.GetBalance(addr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get balance of by user %s. err: %w", user.UserID, err).Error())
	}

	if balance.Int64() < int64(req.Amount) {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("insufficient token balance. require:%d, have:%d", req.Amount, balance.Int64()))
	}

	txHash, err := s.erc20Contract.LockTokens(addr, big.NewInt(int64(req.Amount)), req.Location, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			s.erc20Contract, _ = contract.InitErc20Contract(s.config.PrivateKey)
			txHash, err = s.erc20Contract.LockTokens(addr, big.NewInt(int64(req.Amount)), req.Location, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot bid token at location %s by user %s. new contract instance err: %w", req.Location, user.UserID, err).Error())
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot bid token at location %s by user %s. err: %w", req.Location, user.UserID, err).Error())
		}
	}

	_, err = s.store.CreateBiddingHistory(c.Context(), db.CreateBiddingHistoryParams{
		UserID:    user.UserID,
		Amount:    int64(req.Amount),
		Location:  req.Location,
		TxHash:    txHash.Hex(),
		ExpiresAt: s.nextUnlockDate,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot create bidding history. err: %w", err).Error())
	}

	return c.Status(fiber.StatusOK).JSON(BidTokenResponse{ImageURL: fileURL, Hash: txHash})
}

type SelectedBidderResponse struct {
	UserID string `json:"user_id"`
}

// getSelectedBidder godoc
// @Summary      특정 광고 스팟의 입찰에 성공한 사용자를 불러옵니다.
// @Tags         erc20
// @Param location query string true "광고 스팟"
// @Success      200 {object} api.HighestBidderResponse "사용자와 입찰금액"
// @Router       /selectedBidder [get]
func (s *Server) getSelectedBidder(c *fiber.Ctx) error {
	r := new(GetBidderByLocRequest)
	err := c.QueryParser(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s", err.Error()))
	}
	if errs := utilities.ValidateStruct(r); errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("validation err : %s", errs.Error()))
	}
	addr, _, err := s.erc20Contract.GetCurrentAdOwner(r.Location)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get current AD owner in contract. err: %s", err.Error()))
	}

	user, err := s.store.GetUserByMetamaskAddress(c.Context(), sql.NullString{String: addr.Hex(), Valid: true})

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("cannot find ad owner in this location. err: %s", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("cannot get ad owner in db. err: %s", err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(SelectedBidderResponse{UserID: user.UserID})
}
