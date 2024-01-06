package util

import (
	"github.com/google/uuid"
)

// KeyGen missing godoc.
type KeyGen struct {
}

// NewKeyGen missing godoc.
func NewKeyGen() *KeyGen {
	return &KeyGen{}
}

// Generate missing godoc.
func (k *KeyGen) Generate() string {

	return uuid.NewString()
}
