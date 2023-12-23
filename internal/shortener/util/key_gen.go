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
const KeyLength = 6

func (k *KeyGen) Generate() string {

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	shortKey := make([]byte, KeyLength)
	for i := range shortKey {
		shortKey[i] = Charset[r.Intn(len(Charset))]
	}
	return string(shortKey)
}
