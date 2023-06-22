package utilities

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomFloat(t *testing.T) {
	t.Log(MakeRanFloat(20000, 30000))
}

func TestRandomInt(t *testing.T) {
	t.Log(MakeRanInt(10, 100))
}

func TestMakeRanKey(t *testing.T) {
	var sb strings.Builder
	alphaNum := "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 32; i++ {
		c := alphaNum[MakeRanInt(0, len(alphaNum))]
		sb.WriteByte(c)
	}
	t.Log(sb.String())
	require.Equal(t, 32, len([]byte(sb.String())))
}
