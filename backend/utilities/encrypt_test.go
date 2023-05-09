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
	Start      int64
	Candles    int
	EntryPrice float64
}

func TestEncrypt(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour).UnixMilli(),
		Candles:    500,
		EntryPrice: 5381.432,
	}

	pub, prv := GenerateAsymKey()

	encrypted := EncryptByAsym[interface{}](test, pub)
	t.Log(len(string(encrypted)))

	decrypted := DecryptByAsym(encrypted, prv)
	result := new(testStruct)
	gob.NewDecoder(bytes.NewBuffer(decrypted)).Decode(&result)

	require.Equal(t, test.Start, result.Start)
	require.Equal(t, test.Candles, result.Candles)
	require.Equal(t, test.EntryPrice, result.EntryPrice)
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

func TestMarshalKey(t *testing.T) {
	_, pvk := GenerateAsymKey()
	b := PrivateKeyToBytes(pvk)
	t.Log(b)

	pvk2 := BytesToPrivateKey(b)

	require.Equal(t, pvk, pvk2)
}

func TestGenerateSymKey(t *testing.T) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		t.Error(err)
	}
	t.Log(fmt.Printf("%s", key))
}

func TestAseEncrypt(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour).UnixMilli(),
		Candles:    500,
		EntryPrice: 5381.432,
	}
	jsonData, err := json.Marshal(test)
	require.NoError(t, err)

	encoded := EncrtpByASE(jsonData)

	decodedByte := DecryptByASE(encoded)

	var result testStruct
	err = json.Unmarshal(decodedByte, &result)
	require.NoError(t, err)

	require.Equal(t, test.Start, result.Start)
	require.Equal(t, test.Candles, result.Candles)
	require.Equal(t, test.EntryPrice, result.EntryPrice)
}
