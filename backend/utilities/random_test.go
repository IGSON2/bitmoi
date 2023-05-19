package utilities

import "testing"

func TestRandomFloat(t *testing.T) {
	t.Log(MakeRanFloat(20000, 30000))
}
