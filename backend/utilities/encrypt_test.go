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
	identifier := "ALYJ/z8Bnb4k3SwlZlSr1KAcxn/Km0IYrTKE3fayRnKvKCE2V4BiXe+el4m6g0j2QnBG13z8j0rR5m300pf4CoruyccVqkubqM0KDQOl1S1ZlOw0oSrNi0c2fgMDYREDUzK/Y2eliJwloIsu90GNkky/0ZlVt4usNQM2SBY="

	b := DecryptByASE(identifier)
	json.Unmarshal(b, &c)

	require.NotEmpty(t, c)
}
