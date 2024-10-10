package account

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomID() (string, error) {
	var result strings.Builder
	length := 15

	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHARSET))))
		if err != nil {
			return "", err
		}
		result.WriteByte(CHARSET[index.Int64()])
	}

	return result.String(), nil
}
