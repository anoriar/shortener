package response

// AddURLBatchResponseDTO Результат сохранения нескольких URL пачкой
type AddURLBatchResponseDTO struct {
	CorrelationID string `json:"correlation_id"` // сопоставляет запрос с ответом. Берется из AddURLBatchRequestDTO
	ShortURL      string `json:"short_url"`      // готовая для перехода версия URL http://localhost:8080/s1uHaW"
}
