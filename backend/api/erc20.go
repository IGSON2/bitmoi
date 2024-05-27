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
	FreeAmt         = big.NewInt(3)
	timeoutDuration = 1440
)

type TransactionResponse struct {
	Hash *common.Hash `json:"hash"`
}

// sendFreeErc20 godoc
// @Summary		프리 MOI 토큰을 지급합니다.
// @Tags		erc20
// @Accept		json
// @Produce		json
// @Param		UpdateMetamaskRequest	body		api.UpdateMetamaskRequest	true	"metamask 주소"
// @param 		Authorization header string true "Authorization"
// @Success		200		{object}	string
// @Router      /freeToken [post]
func (s *Server) sendFreeErc20(c *fiber.Ctx) error {
	r := new(UpdateMetamaskRequest)
	err := c.BodyParser(r)
	if errs := utilities.ValidateStruct(r); err != nil || errs != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("parsing err : %s, validation err : %s", err, errs.Error()))
	}

	payload, ok := c.Locals(authorizationPayloadKey).(*token.Payload)
	if !ok {
		return fmt.Errorf("cannot get authorization payload")
	}
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

	if timeoutUnix := s.faucetTimeouts[user.UserID]; time.Now().After(time.Unix(timeoutUnix, 0)) {
		ToAddr := common.HexToAddress(user.MetamaskAddress.String)

		hash, err := s.erc20Contract.SendFreeTokens(ToAddr, FreeAmt, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
		if err != nil {
			if strings.Contains(err.Error(), "nonce") {
				s.erc20Contract, err = contract.InitErc20Contract(s.config.PrivateKey)
				if err != nil {
					log.Error().Err(err).Msg("cannot regenerate erc20 contract instance.")
					return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
				}
			}
			hash, err = s.erc20Contract.SendFreeTokens(ToAddr, FreeAmt, contract.TransactOptions{GasLimit: contract.DefaultGasLimit})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}

		duration := time.Duration(timeoutDuration) * time.Minute
		s.faucetTimeouts[user.UserID] = time.Now().Add(duration).Unix()
		return c.Status(fiber.StatusOK).JSON(TransactionResponse{hash})
	} else {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("%s left until next allowance", common.PrettyDuration(time.Until(time.Unix(timeoutUnix, 0)))).Error())
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
			erc20, reGenErc20Err := contract.InitErc20Contract(s.config.PrivateKey)
			if reGenErc20Err != nil {
				return nil, fmt.Errorf("cannot init erc20 contract : %w", reGenErc20Err)
			}
			log.Warn().Msgf("Erc20 instance was regenerated caused by this err: %s", err.Error())
			s.erc20Contract = erc20
		}
		return nil, fmt.Errorf("cannot update token ledger. err: %s", err.Error())
	}
	return txResult.TxHash, nil
}
