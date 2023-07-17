package evmcleint

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

func NewEvmClient() {
	rpcAddress := "https://api.baobab.klaytn.net:8651"

	c, err := ethclient.Dial(rpcAddress)
	if err != nil {
		log.Err(err).Msg("cannot dial to rpc node")
		return
	}
	ctx := context.Background()
	chainID, err := c.ChainID(ctx)
	if err != nil {
		log.Err(err).Msg("cannot get chain id")
	}
	fmt.Println(chainID, chainID.Int64())
}
