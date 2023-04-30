package utilities

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/gob"
)

func GenerateKey() (*rsa.PublicKey, *rsa.PrivateKey) {
	pvk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil
	}
	return &pvk.PublicKey, pvk
}

func Encrypt[T any](data T, pub *rsa.PublicKey) []byte {
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

func Decrypt(b []byte, priv *rsa.PrivateKey) []byte {
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
