package request

type AddURLBatchRequestDTO struct {
	AuthRequest
	Items []AddURLBatchItemDTO
}

type AddURLBatchItemDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
