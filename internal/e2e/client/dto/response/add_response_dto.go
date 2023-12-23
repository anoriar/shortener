package response

// AddResponseDto missing godoc.
type AddResponseDto struct {
	BaseResponseDto
	Body string
}

// NewShortenerResponseDto missing godoc.
func NewShortenerResponseDto(statusCode int, contentType string, token string, body string) *AddResponseDto {
	return &AddResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
