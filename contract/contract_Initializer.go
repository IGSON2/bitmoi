package contract

import (
	"bitmoi/contract/evmclient"
	"bitmoi/contract/evmgaspricer"
	"bitmoi/contract/signer"
	"bitmoi/contract/transaction"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

const (
	ERC20Address = "0xC6D6e6E081dA15bd4258bDbdC121AD2D72ad793B"
)

var (
	gasPrice = big.NewInt(25000000000)
)

func InitializeClient(
	senderKeyPair *signer.Keypair,
) (*evmclient.EvmClient, error) {
	ethClient, err := evmclient.NewEvmClient(senderKeyPair)
	if err != nil {
		log.Error().Err(fmt.Errorf("eth client initialization error: %v", err))
		return nil, err
	}
	return ethClient, nil
}

func InitializeTransactor(
	gasPrice *big.Int,
	txFabric TxFabric,
	client *evmclient.EvmClient,
) (Transactor, error) {
	var trans Transactor

	gasPricer := evmgaspricer.NewLondonGasPriceClient(
		client,
		&evmgaspricer.GasPricerOpts{UpperLimitFeePerGas: gasPrice},
	)
	trans = NewSignAndSendTransactor(txFabric, gasPricer, client)

	return trans, nil
}

func InitErc20Contract(privKey string) (*ERC20Contract, error) {

	sender, err := signer.NewKeypairFromPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	c, err := InitializeClient(sender)
	if err != nil {
		return nil, err
	}
	t, err := InitializeTransactor(gasPrice, transaction.NewTransaction, c)
	if err != nil {
		return nil, err
	}
	return NewERC20Contract(c, common.HexToAddress(ERC20Address), t), nil
}
