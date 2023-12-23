package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// SignService missing godoc.
type SignService struct {
	SecretKey string
}

// NewSignService missing godoc.
func NewSignService(secretKey string) *SignService {
	return &SignService{SecretKey: secretKey}
}

// Sign missing godoc.
func (ss *SignService) Sign(data []byte) string {
	h := hmac.New(sha256.New, []byte(ss.SecretKey))
	h.Write(data)
	signedData := append(data, h.Sum(nil)...)
	return hex.EncodeToString(signedData)
}

// Verify missing godoc.
func (ss *SignService) Verify(message []byte, signature []byte) bool {
	h := hmac.New(sha256.New, []byte(ss.SecretKey))
	h.Write(message)
	expectedSignature := h.Sum(nil)

	return hmac.Equal(signature, expectedSignature)
}

// Decode missing godoc.
func (ss *SignService) Decode(signedData string) (message []byte, signature []byte, err error) {
	decodedData, err := hex.DecodeString(signedData)
	if err != nil {
		return nil, nil, err
	}
	// Получаем сообщение и сигнатуру из подписанной куки.
	message = decodedData[:len(decodedData)-sha256.Size]
	signature = decodedData[len(decodedData)-sha256.Size:]
	return message, signature, err
}
