package response

type DeleteURLBatchResponseDto struct {
	StatusCode  int
	ContentType string
}

func NewDeleteURLBatchResponseDto(statusCode int, contentType string) *DeleteURLBatchResponseDto {
	return &DeleteURLBatchResponseDto{StatusCode: statusCode, ContentType: contentType}
}
