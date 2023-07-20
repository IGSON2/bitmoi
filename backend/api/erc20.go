package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
)

var (
	FreeAmt         = big.NewInt(1)
	timeoutDuration = 3
)

type TransactionResponse struct {
	Hash *common.Hash `json:"hash"`
}

func (s *Server) sendFreeErc20(c *fiber.Ctx) error {

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
	}

	if !user.MetamaskAddress.Valid {
		return c.Status(fiber.StatusForbidden).SendString("cannot get metamask address by user")
	}

	if timeout := s.erc20Contract.Timeouts[user.MetamaskAddress.String]; time.Now().After(timeout) {
		ToAddr := common.HexToAddress(user.MetamaskAddress.String)

		hash, err := s.erc20Contract.SendFreeTokens(ToAddr, FreeAmt, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		timeout := time.Duration(timeoutDuration) * time.Minute
		grace := timeout / 288

		s.erc20Contract.Timeouts[user.MetamaskAddress.String] = time.Now().Add(timeout - grace)
		return c.Status(fiber.StatusOK).JSON(TransactionResponse{hash})
	} else {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("%s left until next allowance", common.PrettyDuration(time.Until(timeout))).Error())
	}

}

func (s *Server) spendErc20OnComp(c *fiber.Ctx, scoreId string) (*common.Hash, error) {
	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("cannot get authorization payload")
	}

	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return nil, fmt.Errorf("cannot get user by authorization payload")
	}

	fromAddr := common.HexToAddress(user.MetamaskAddress.String)

	balance, _ := s.erc20Contract.GetBalance(fromAddr)
	if balance == nil || balance.Cmp(big.NewInt(1)) < 0 {
		return nil, fmt.Errorf("insufficient balance. Addr: %s, %d MOI", fromAddr.Hex(), balance.Int64())
	}

	hash, err := s.erc20Contract.SpendTokens(fromAddr, FreeAmt, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
	if err != nil {
		return nil, err
	}

	_, err = s.store.CreateUsedToken(c.Context(), db.CreateUsedTokenParams{
		ScoreID:         scoreId,
		UserID:          user.UserID,
		MetamaskAddress: user.MetamaskAddress.String,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update metamask address. err: %s", err.Error())
	}
	return hash, nil
}