package util

import (
	"math/rand"
	"time"
)

type KeyGen struct {
}

func NewKeyGen() *KeyGen {
	return &KeyGen{}
}

const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const KeyLength = 10

func (k *KeyGen) Generate() string {

	rand.Seed(time.Now().UnixNano())
	shortKey := make([]byte, KeyLength)
	for i := range shortKey {
		shortKey[i] = Charset[rand.Intn(len(Charset))]
	}
	return string(shortKey)
}
