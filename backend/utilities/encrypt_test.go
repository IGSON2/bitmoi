package utilities

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Start      int64   `json:"start"`
	Candles    int     `json:"candles"`
	EntryPrice float64 `json:"entryprice"`
	Secret     string
}

func TestBase64(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour).UnixMilli(),
		Candles:    500,
		EntryPrice: 5381.432,
	}
	nb := bytes.NewBuffer(nil)
	gob.NewEncoder(nb).Encode(test)

	t.Log(len(Base64Encode(nb.Bytes())))
}

func TestGenerateSymKey(t *testing.T) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}
	t.Log(fmt.Printf("%s", key))
}

func TestAESEncrypt(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour).UnixMilli(),
		Candles:    500,
		EntryPrice: 5381.432,
		Secret:     "s",
	}
	encoded := EncrtpByASE(test)

	decodedByte := DecryptByASE(encoded)

	var result testStruct
	err := json.Unmarshal(decodedByte, &result)
	require.NoError(t, err)

	require.Equal(t, test.Start, result.Start)
	require.Equal(t, test.Candles, result.Candles)
	require.Equal(t, test.EntryPrice, result.EntryPrice)
	require.Equal(t, test.Secret, result.Secret, fmt.Sprintf("origin : %v, result :%v", test.Secret, result.Start))
}

func TestAESDecrypt(t *testing.T) {
	var c IdentificationData
	identifier := "00b609ff3f019dbe24c52f266654abd4a01cc67fca9b4218ad3284ddf6b24372af28212653857d47f69e948ef9d409e2407247db7ff28f4ec5fd2daa9984ef0588e8dbc904fc1ed9a4c7080803bbd027599bf920e57594ca1861204e4263111d0465ea2924b49d9d27e3c26abd10c9c65fb1d19e51b181a8330933495ce937e2153a599da4fc1a35e7a0754c0ceba185c68e4cdd6a0738d6c65b"

	b := DecryptByASE(identifier)
	json.Unmarshal(b, &c)

	require.NotEmpty(t, c)
}
