package api

import (
	"bitmoi/backend/contract"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	ErrClosedBiddingLoop = errors.New("bidding loop has been closed")
)

type NextUnlockResponse struct {
	NextUnlock int64 `json:"next_unlock"`
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
				continue
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
		NextUnlock: s.nextUnlockDate.Unix()})
}
