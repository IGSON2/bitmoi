package utilities

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type IdentificationData struct {
	Name         string  `json:"name"`
	Interval     string  `json:"interval"`
	RefTimestamp int64   `json:"reftimestamp"`
	PriceFactor  float64 `json:"pricefactor"`
	VolumeFactor float64 `json:"volumefactor"`
	TimeFactor   int64   `json:"ranpastdate"`
}

func (i *IdentificationData) IsZeroFactor() bool {
	return i.PriceFactor == 0 || i.TimeFactor == 0 || i.VolumeFactor == 0
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

	return Base64Encode(encryptedData)
}

func DecryptByASE(encrypted string) []byte {
	b := Base64Decode(encrypted)
	block, _ := aes.NewCipher([]byte(GetConfig("../../.").SymmetricKey))
	iv := make([]byte, aes.BlockSize)

	stream := cipher.NewCTR(block, iv)
	decryptedData := make([]byte, len(b))

	stream.XORKeyStream(decryptedData, b)
	return decryptedData
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("faild to hash password! err : %v", err)
	}
	return string(hashed), nil
}

func CheckPassword(password, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
