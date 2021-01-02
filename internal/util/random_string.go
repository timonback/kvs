package util

import (
	"crypto/rand"
	"math/big"
)

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		bigInt, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		s[i] = letters[bigInt.Int64()]
	}
	return string(s)
}
