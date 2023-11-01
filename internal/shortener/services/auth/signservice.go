package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type SignService struct {
	SecretKey string
}

func NewSignService(secretKey string) *SignService {
	return &SignService{SecretKey: secretKey}
}

func (ss *SignService) Sign(data []byte) string {
	h := hmac.New(sha256.New, []byte(ss.SecretKey))
	h.Write(data)
	signedData := append(data, h.Sum(nil)...)
	return hex.EncodeToString(signedData)
}

func (ss *SignService) Verify(message []byte, signature []byte) bool {
	h := hmac.New(sha256.New, []byte(ss.SecretKey))
	h.Write(message)
	expectedSignature := h.Sum(nil)

	return hmac.Equal(signature, expectedSignature)
}

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
