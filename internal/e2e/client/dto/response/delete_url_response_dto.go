package response

// DeleteURLBatchResponseDto missing godoc.
type DeleteURLBatchResponseDto struct {
	BaseResponseDto
}

// NewDeleteURLBatchResponseDto missing godoc.
func NewDeleteURLBatchResponseDto(statusCode int, contentType string, token string) *DeleteURLBatchResponseDto {
	return &DeleteURLBatchResponseDto{BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}}
}
