package token

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

func GenerateRecCode() (string, error) {
	b, err := uuid.New().MarshalBinary()
	if err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(b[:4])), nil
}
