package util

import (
	"math/rand"
	"time"
)

const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const KeyLength = 6

func GenerateShortKey() string {

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, KeyLength)
	for i := range shortKey {
		shortKey[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(shortKey)
}
