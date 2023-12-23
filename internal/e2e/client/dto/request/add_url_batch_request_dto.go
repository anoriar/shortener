package request

// AddURLBatchRequestDTO missing godoc.
type AddURLBatchRequestDTO struct {
	AuthRequest
	Items []AddURLBatchItemDTO
}

// AddURLBatchItemDTO missing godoc.
type AddURLBatchItemDTO struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
