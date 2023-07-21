package contract

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPrivateKey = "2b6399ecbebd49c4fdf3ddfab539b306f1f007dfeb1dd1e9fdcd375ccc04f788"
)

func newTestErc20Contract(t *testing.T) *ERC20Contract {
	contract, err := InitErc20Contract(testPrivateKey)
	require.NoError(t, err)
	return contract
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
