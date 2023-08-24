package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"database/sql"
	"errors"
	"fmt"
	"math/big"
	"os"
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

type HighestBidderResponse struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}

func (s *Server) BiddingLoop() error {

	for {
		biddingTimer := time.NewTimer(s.config.BiddingDuration)
		select {
		case <-biddingTimer.C:
			hash, err := s.erc20Contract.UnLockTokens(contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
			if err != nil {
				var initErr error
				s.erc20Contract, initErr = contract.InitErc20Contract(s.config.PrivateKey)
				if initErr != nil {
					log.Err(err).Msgf("Cannot unlock token. stop server..")
					os.Exit(100)
				}
			}
			receipt, err := s.erc20Contract.WaitAndReturnTxReceipt(hash)
			if err != nil || receipt == nil {
				var initErr error
				s.erc20Contract, initErr = contract.InitErc20Contract(s.config.PrivateKey)
				if initErr != nil {
					log.Err(err).Msgf("Cannot get receipt of unlock token transaction. stop server..")
					os.Exit(100)
				}
				continue
			}
			s.nextUnlockDate = time.Now().Add(s.config.BiddingDuration)
			log.Info().Msgf("Unlock token successfully. hash: %s, next unlock date: %s", hash.Hex(), s.nextUnlockDate.Format("2006-01-02 15:04:05"))
		case <-s.exitCh:
			return ErrClosedBiddingLoop
		}
	}
}

// getNextUnlockDate godoc
// @Summary      Get next unlock date
// @Description  Get next date of unlock ad bidding
// @Tags         erc20
// @Success      200 {object} api.NextUnlockResponse "unix timestamp"
// @Router       /nextBidUnlock [get]
func (s *Server) getNextUnlockDate(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(NextUnlockResponse{
		NextUnlock: s.nextUnlockDate.Format("2006-01-02T15:04:05")})
}

// getHighestBidder godoc
// @Summary      Get highest bidder
// @Description  Get highest bidder by specific location
// @Tags         erc20
// @Param location query string true "Location for the spot to advertise"
// @Success      200 {object} api.HighestBidderResponse "bidder and amount of tokens bid"
// @Router       /highestBidder [get]
func (s *Server) getHighestBidder(c *fiber.Ctx) error {
	req := new(GetHighestBidderRequest)
	err := c.QueryParser(req)
	if errs := utilities.ValidateStruct(req); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}
	addr, amt, err := s.erc20Contract.GetHighestBidder(req.Location)
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

// bidToken godoc
//
//	@Summary		Bid token
//	@Description	Bid token to specified location
//	@Tags			erc20
//	@Accept			json
//	@Produce		json
//	@Param			request	body		api.BidTokenRequest	true	"bidding request"
//	@param Authorization header string true "Authorization"
//	@Success		200		{object}	api.ScoreResponse
//	@Router       /bidToken [post]
func (s *Server) bidToken(c *fiber.Ctx) error {
	req := new(BidTokenRequest)
	err := c.BodyParser(req)
	if errs := utilities.ValidateStruct(req); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
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
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot bid token at location %s by user %s. err: %w", req.Location, user.UserID, err).Error())
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

	return c.Status(fiber.StatusOK).JSON(TransactionResponse{Hash: txHash})
}
