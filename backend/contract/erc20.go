package contract

import (
	"bitmoi/backend/contract/consts"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type ERC20Contract struct {
	Contract
	Timeouts map[string]time.Time
}

func NewERC20Contract(
	client ContractCallerDispatcher,
	erc20ContractAddress common.Address,
	transactor Transactor,
) *ERC20Contract {
	a, _ := abi.JSON(strings.NewReader(consts.ERC20ABI))
	b := common.FromHex(consts.ERC20Bin)
	return &ERC20Contract{
		Contract: NewContract(erc20ContractAddress, a, b, client, transactor),
		Timeouts: make(map[string]time.Time),
	}
}

func (c *ERC20Contract) WaitAndReturnTxReceipt(hash *common.Hash) (*types.Receipt, error) {
	return c.Contract.client.WaitAndReturnTxReceipt(*hash)
}

func (c *ERC20Contract) GetBalance(address common.Address) (*big.Int, error) {
	log.Debug().Msgf("Getting balance for %s", address.String())
	res, err := c.CallContract("balanceOf", address)
	if err != nil {
		return nil, err
	}
	b := abi.ConvertType(res[0], new(big.Int)).(*big.Int)
	return b, nil
}

func (c *ERC20Contract) SpendTokens(
	from common.Address,
	amount *big.Int,
	opts TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Spend %s tokens from %s", amount.String(), from.String())
	return c.ExecuteTransaction("spendToken", opts, from, amount)
}

func (c *ERC20Contract) SendFreeTokens(
	to common.Address,
	amount *big.Int,
	opts TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Sending %s tokens to %s", amount.String(), to.String())
	return c.ExecuteTransaction("sendFreeToken", opts, to, amount)
}

func (c *ERC20Contract) LockTokens(
	from common.Address,
	amount *big.Int,
	adLocation string,
	opts TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Locked %s tokens from %s for bidding location %s", amount.String(), from.String(), adLocation)
	return c.ExecuteTransaction("lockTokens", opts, from, amount, adLocation)
}

func (c *ERC20Contract) UnLockTokens(
	opts TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("UnLock all tokens")
	return c.ExecuteTransaction("unlockTokens", opts)
}
