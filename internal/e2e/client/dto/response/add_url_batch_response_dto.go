package response

type AddURLBatchItemDTO struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type AddURLBatchResponseDto struct {
	BaseResponseDto
	Body []AddURLBatchItemDTO
}

func NewAddURLBatchResponseDto(statusCode int, contentType string, token string, body []AddURLBatchItemDTO) *AddURLBatchResponseDto {
	return &AddURLBatchResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
