package response

// DeleteUserURLsResponseDto missing godoc.
type DeleteUserURLsResponseDto struct {
	BaseResponseDto
}

// NewDeleteUserURLsResponseDto missing godoc.
func NewDeleteUserURLsResponseDto(statusCode int, contentType string, token string) *DeleteUserURLsResponseDto {
	return &DeleteUserURLsResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode: statusCode,
		Token:      token,
	}}
}
