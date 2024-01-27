package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenRecCode(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	for i := 0; i < 10; i++ {
		b, err := GenerateRecCode()
		require.NoError(t, err)

		t.Log(b)
	}
}
