package util

// KeyGenInterface missing godoc.
//
//go:generate mockgen -source=key_gen_interface.go -destination=mock/key_gen.go -package=mock KeyGenInterface
type KeyGenInterface interface {
	Generate() string
}
