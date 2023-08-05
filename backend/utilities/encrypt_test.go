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
	if testing.Short() {
		t.Skip()
	}
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}
}

func TestAESEncrypt(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
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

func TestEncodeIdentifier(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	identifier := "gy+Itl5buecRyGNVDCUEHZ/iRz8DpEQDnaXypZskcqc0JxkoPkltkTPNQa/G/Yv6uf86UtMwBBMHdmMdWyS5z+UjTqaYF0gY7bXn7J9ITiNF62lOgCtCjlM8aw56KiK57S54TIHEAYOSiL2eEv35zVwHk38ol81ECQMH2AjxLxtNHiQHr3CsaNp/qg=="

	info := new(IdentificationData)

	require.NoError(t, json.Unmarshal(DecryptByASE(identifier), info))

	fmt.Println(info)

}
