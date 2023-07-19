package evmcleint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPrivateKey = "2b6399ecbebd49c4fdf3ddfab539b306f1f007dfeb1dd1e9fdcd375ccc04f788"
)

func TestNewClient(t *testing.T) {
	keyPair, err := NewKeypairFromPrivateKey(testPrivateKey)
	require.NoError(t, err)
	client := NewEvmClient(keyPair)
	require.Equal(t, int64(1001), client.chainId.Int64())
	t.Log(client.chainId.Int64())
}
