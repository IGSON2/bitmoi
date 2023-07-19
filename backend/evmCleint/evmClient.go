package evmcleint

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
)

const (
	TestRpcURL = "https://api.baobab.klaytn.net:8651"
)

type Signer interface {
	CommonAddress() common.Address

	// Sign calculates an ECDSA signature.
	// The produced signature must be in the [R || S || V] format where V is 0 or 1.
	Sign(digestHash []byte) ([]byte, error)
}

type EvmClient struct {
	rpcClient *rpc.Client
	chainId   *big.Int
	signer    Signer
	nonce     *big.Int
	nonceLock sync.Mutex
}

func NewEvmClient(s Signer) *EvmClient {
	ctx := context.Background()
	rpcClient, err := rpc.DialContext(ctx, TestRpcURL)
	if err != nil {
		log.Err(err).Msg("cannot dial to rpc node")
		return nil
	}

	chainId := new(hexutil.Big)
	err = rpcClient.CallContext(ctx, chainId, "eth_chainId")
	if err != nil {
		log.Err(err).Msg("cannot get chain id")
		return nil
	}

	client := EvmClient{
		rpcClient: rpcClient,
		chainId:   (*big.Int)(chainId),
		signer:    s,
	}

	return &client
}

// TODO: Contract 객체 생성 후 Contract 메서드 호출해야 함
func (c *EvmClient) SendToken(address common.Address) (*common.Hash, error) {
	// c.LockNonce()
	// n, err := c.UnsafeNonce()
	// if err != nil {
	// 	c.UnlockNonce()
	// 	return &common.Hash{}, err
	// }

	// types.NewTransaction(n, address)
	// types.SignTx()
	// eip155Signer := types.NewEIP155Signer(c.chainId)

	// c.UnlockNonce()
	return nil, nil
}

func (c *EvmClient) LockNonce() {
	c.nonceLock.Lock()
}

func (c *EvmClient) UnlockNonce() {
	c.nonceLock.Unlock()
}

func (c *EvmClient) UnsafeNonce() (*big.Int, error) {
	var err error
	for i := 0; i <= 10; i++ {
		if c.nonce == nil {
			nonce, err := c.PendingNonceAt(context.Background(), c.signer.CommonAddress())
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			c.nonce = big.NewInt(0).SetUint64(nonce)
			return c.nonce, nil
		}
		return c.nonce, nil
	}
	return nil, err
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (c *EvmClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := c.rpcClient.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}
