package utilities

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func MakeRanInt(minimum, maximum int) int {
	ranSeed := big.NewInt(int64(maximum - minimum))
	if ranSeed.Cmp(big.NewInt(0)) <= 0 {
		log.Error().Msgf("Random int seed is equal or smaller than 0. seed: %d", maximum-minimum)
		return 0
	}
	ranBigNum, err := rand.Int(rand.Reader, ranSeed)
	if err != nil {
		return -1
	}
	return int(ranBigNum.Int64()) + minimum
}

func MakeRanFloat(minimum, maximum int) (float64, error) {
	if maximum <= minimum {
		return 0, fmt.Errorf("maximum number is equal smaller than minimum maxmimum: %d, minimum: %d", maximum, minimum)
	}

	ranSeed := big.NewInt(int64(maximum - minimum))
	ranFloat := mathrand.Float64()

	return float64(ranSeed.Int64())*ranFloat + float64(minimum), nil
}

func MakeRanString(length int) string {
	var sb strings.Builder

	for i := 0; i < length; i++ {
		c := alphabet[MakeRanInt(0, len(alphabet))]
		sb.WriteByte(c)
	}

	return sb.String()
}

func MakeRanEmail() string {
	return fmt.Sprintf("%s@email.com", MakeRanString(6))
}
