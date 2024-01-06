package response

// GetResponseDto missing godoc.
type GetResponseDto struct {
	BaseResponseDto
	Location string
}

// NewGetResponseDto missing godoc.
func NewGetResponseDto(statusCode int, contentType string, location string, token string) *GetResponseDto {
	return &GetResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode: statusCode,
		Token:      token,
	}, Location: location}
}
