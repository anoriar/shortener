package util

import (
	"math/rand"
	"time"
)

// KeyGen missing godoc.
type KeyGen struct {
}

// NewKeyGen missing godoc.
func NewKeyGen() *KeyGen {
	return &KeyGen{}
}

// Charset missing godoc.
const Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// KeyLength missing godoc.
const KeyLength = 6

// Generate missing godoc.
func (k *KeyGen) Generate() string {

	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	shortKey := make([]byte, KeyLength)
	for i := range shortKey {
		shortKey[i] = Charset[r.Intn(len(Charset))]
	}
	return string(shortKey)
}
