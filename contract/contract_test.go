package contract

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestContract(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	contract, err := InitErc20Contract(testPrivateKey)
	require.NoError(t, err)
	require.NotNil(t, contract)

	var lastHash *common.Hash

	addrArray := testGenerateAddress(t, 10)
	for i, addr := range addrArray {
		hash, err := contract.SendFreeTokens(addr, big.NewInt(1), TransactOptions{GasLimit: DefaultGasLimit})
		require.NoError(t, err)
		require.NotNil(t, hash)
		if i == len(addrArray)-1 {
			lastHash = hash
		}
	}

	_, err = contract.client.WaitAndReturnTxReceipt(*lastHash)
	require.NoError(t, err)

	for _, addr := range addrArray {
		balance, err := contract.GetBalance(addr)
		require.NoError(t, err)
		require.GreaterOrEqual(t, balance.Int64(), int64(1))
	}

}

func testGenerateAddress(t *testing.T, cnt int) []common.Address {
	addrArray := []common.Address{}
	for i := 0; i < cnt; i++ {
		pv, err := crypto.GenerateKey()
		require.NoError(t, err)
		pub := pv.Public().(*ecdsa.PublicKey)
		addr := crypto.PubkeyToAddress(*pub)
		addrArray = append(addrArray, addr)
	}
	return addrArray
}
