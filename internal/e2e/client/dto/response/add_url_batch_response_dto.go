package response

// AddURLBatchItemDTO missing godoc.
type AddURLBatchItemDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// AddURLBatchResponseDto missing godoc.
type AddURLBatchResponseDto struct {
	BaseResponseDto
	Body []AddURLBatchItemDTO
}

// NewAddURLBatchResponseDto missing godoc.
func NewAddURLBatchResponseDto(statusCode int, contentType string, token string, body []AddURLBatchItemDTO) *AddURLBatchResponseDto {
	return &AddURLBatchResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
