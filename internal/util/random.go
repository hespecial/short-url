package util

import (
	"math/rand"
	"time"
)

const (
	seed = "1234567890qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
)

var (
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func GenerateRandomBytes(minLen, maxLen int) []byte {
	length := rand.Intn(maxLen-minLen+1) + minLen
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = seed[r.Intn(len(seed))]
	}
	return result
}
