package response

type AddResponseDto struct {
	BaseResponseDto
	Body string
}

func NewShortenerResponseDto(statusCode int, contentType string, token string, body string) *AddResponseDto {
	return &AddResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
