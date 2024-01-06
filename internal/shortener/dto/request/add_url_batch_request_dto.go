package request

// AddURLBatchRequestDTO Запрос на сохранение нескольких URL
type AddURLBatchRequestDTO struct {
	CorrelationID string `json:"correlation_id"` // сопоставляет результаты ответа с запросом. Не хранится
	OriginalURL   string `json:"original_url"`   // оригинальный URL для сохранения
}
