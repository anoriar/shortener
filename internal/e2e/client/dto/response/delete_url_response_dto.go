package response

type DeleteURLBatchResponseDto struct {
	BaseResponseDto
}

func NewDeleteURLBatchResponseDto(statusCode int, contentType string, token string) *DeleteURLBatchResponseDto {
	return &DeleteURLBatchResponseDto{BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}}
}
