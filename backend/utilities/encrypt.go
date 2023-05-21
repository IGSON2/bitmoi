package utilities

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
)

type IdentificationData struct {
	Name         string  `json:"name"`
	Interval     string  `json:"interval"`
	RefTimestamp int64   `json:"reftimestamp"`
	PriceFactor  float64 `json:"pricefactor"`
	VolumeFactor float64 `json:"volumefactor"`
	TimeFactor   int64   `json:"ranpastdate"`
}

func GenerateAsymKey() (*rsa.PublicKey, *rsa.PrivateKey) {
	pvk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil
	}
	return &pvk.PublicKey, pvk
}

func EncryptByAsym[T any](data T, pub *rsa.PublicKey) []byte {
	buffer := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buffer).Encode(data)
	if err != nil {
		return nil
	}
	b, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, buffer.Bytes(), nil)
	if err != nil {
		return nil
	}
	return b
}

func DecryptByAsym(b []byte, priv *rsa.PrivateKey) []byte {
	r, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, b, nil)
	if err != nil {
		return nil
	}
	return r
}

func Base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base64Decode(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}

func PrivateKeyToBytes(prv *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(prv)
}

func BytesToPrivateKey(b []byte) *rsa.PrivateKey {
	prv, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil
	}
	return prv
}

func EncrtpByASE[T any](data T) string {
	bytesData, err := json.Marshal(data)
	if err != nil {
		panic(err.Error())
	}

	block, err := aes.NewCipher([]byte(GetConfig("../../.").SymmetricKey))
	if err != nil {
		panic(err.Error())
	}
	iv := make([]byte, aes.BlockSize)
	stream := cipher.NewCTR(block, iv)

	encryptedData := make([]byte, len(bytesData))
	stream.XORKeyStream(encryptedData, bytesData)

	return hex.EncodeToString(encryptedData)
}

func DecryptByASE(encrypted string) []byte {
	b, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil
	}
	block, _ := aes.NewCipher([]byte(GetConfig("../../.").SymmetricKey))
	iv := make([]byte, aes.BlockSize)

	stream := cipher.NewCTR(block, iv)
	decryptedData := make([]byte, len(b))

	stream.XORKeyStream(decryptedData, b)
	return decryptedData
}
