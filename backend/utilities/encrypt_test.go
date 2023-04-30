package utilities

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Start      time.Time
	Candles    int
	EntryPrice float64
}

func TestEncrypt(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour),
		Candles:    500,
		EntryPrice: 5381.432,
	}

	pub, prv := GenerateKey()

	encrypted := Encrypt[interface{}](test, pub)
	t.Log(len(string(encrypted)))

	decrypted := Decrypt(encrypted, prv)
	result := new(testStruct)
	gob.NewDecoder(bytes.NewBuffer(decrypted)).Decode(&result)

	require.WithinDuration(t, test.Start, result.Start, 1*time.Second)
	require.Equal(t, test.Candles, result.Candles)
	require.Equal(t, test.EntryPrice, result.EntryPrice)
}

func TestBase64(t *testing.T) {
	test := testStruct{
		Start:      time.Now().Add(-100 * time.Hour),
		Candles:    500,
		EntryPrice: 5381.432,
	}
	nb := bytes.NewBuffer(nil)
	gob.NewEncoder(nb).Encode(test)

	t.Log(len(Base64Encode(nb.Bytes())))
}

func TestMarshalKey(t *testing.T) {
	_, pvk := GenerateKey()
	b := PrivateKeyToBytes(pvk)
	t.Log(b)

	pvk2 := BytesToPrivateKey(b)

	require.Equal(t, pvk, pvk2)
}
