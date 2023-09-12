package dto

type AddResponseDto struct {
	StatusCode  int
	ContentType string
	Body        string
}

func NewShortenerResponseDto(statusCode int, contentType string, body string) *AddResponseDto {
	return &AddResponseDto{StatusCode: statusCode, ContentType: contentType, Body: body}
}
