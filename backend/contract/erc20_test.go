package contract

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestLockToken(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	contract := newTestErc20Contract(t)

	addrString := "0xB2851696045E2097C6abb8af074eee432e42aEf7"
	fromAddr := common.HexToAddress(addrString)
	hash, err := contract.LockTokens(fromAddr, big.NewInt(5), "practice", TransactOptions{GasLimit: DefaultGasLimit})

	require.NoError(t, err)
	require.NotNil(t, hash)

	_, err = contract.client.WaitAndReturnTxReceipt(*hash)
	require.NoError(t, err)

	// unlockHash, err := contract.UnLockTokens(TransactOptions{GasLimit: DefaultGasLimit})
	// require.NoError(t, err)
	// require.NotNil(t, unlockHash)

	// _, err = contract.client.WaitAndReturnTxReceipt(*unlockHash)
	// require.NoError(t, err)
}

func TestGetHighestBidder(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	contract := newTestErc20Contract(t)

	addr, amt, err := contract.GetHighestBidder("practice")
	require.NoError(t, err)
	fmt.Println(addr.Hex(), amt.Int64())
}
