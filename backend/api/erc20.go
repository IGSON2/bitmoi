package api

import (
	"bitmoi/backend/contract"
	db "bitmoi/backend/db/sqlc"
	"bitmoi/backend/token"
	"bitmoi/backend/utilities"
	"database/sql"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	FreeAmt         = big.NewInt(1)
	timeoutDuration = 1440
)

type TransactionResponse struct {
	Hash *common.Hash `json:"hash"`
}

// sendFreeErc20 godoc
//
//	@Summary		Post sendFreeErc20
//	@Description	request for free token
//	@Tags			erc20
//	@Accept			json
//	@Produce		json
//	@Param			MetamaskAddressRequest	body		api.MetamaskAddressRequest	true	"eth address"
//	@param Authorization header string true "Authorization"
//	@Success		200		{object}	api.TransactionResponse
//	@Router       /freetoken [post]
func (s *Server) sendFreeErc20(c *fiber.Ctx) error {
	r := new(MetamaskAddressRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload := c.Locals(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUser(c.Context(), payload.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot get user by token payload. err: %w", err).Error())
	}

	if !user.MetamaskAddress.Valid {
		log.Info().Msgf("%s dosen't have account, init new account.", user.UserID)
		_, err = s.store.UpdateUserMetamaskAddress(c.Context(), db.UpdateUserMetamaskAddressParams{
			MetamaskAddress:  sql.NullString{String: r.Addr, Valid: true},
			UserID:           payload.UserID,
			AddressChangedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Errorf("cannot update user address. err: %w", err).Error())
		}
		user.MetamaskAddress.String = r.Addr
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

	spendTxArg := db.SpendTokenTxParams{
		CreateUsedTokenParams: db.CreateUsedTokenParams{
			ScoreID:         scoreId,
			UserID:          user.UserID,
			MetamaskAddress: user.MetamaskAddress.String,
		},
		Contract: s.erc20Contract,
		FromAddr: fromAddr,
		Amount:   FreeAmt,
	}
	txResult, err := s.store.SpendTokenTx(c.Context(), spendTxArg)

	if err != nil {
		if strings.Contains(err.Error(), "nonce") {
			//TODO
		}
		return nil, fmt.Errorf("cannot update token ledger. err: %s", err.Error())
	}
	return txResult.TxHash, nil
}
