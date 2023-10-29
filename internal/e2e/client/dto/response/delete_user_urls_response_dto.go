package response

type DeleteUserURLsResponseDto struct {
	BaseResponseDto
}

func NewDeleteUserURLsResponseDto(statusCode int, contentType string, token string) *DeleteUserURLsResponseDto {
	return &DeleteUserURLsResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode: statusCode,
		Token:      token,
	}}
}
