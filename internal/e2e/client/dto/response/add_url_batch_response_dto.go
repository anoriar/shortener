package response

type AddURLBatchItemDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type AddURLBatchResponseDto struct {
	StatusCode  int
	ContentType string
	Body        []AddURLBatchItemDTO
}

func NewAddURLBatchResponseDto(statusCode int, contentType string, body []AddURLBatchItemDTO) *AddURLBatchResponseDto {
	return &AddURLBatchResponseDto{StatusCode: statusCode, ContentType: contentType, Body: body}
}
