package response

type GetResponseDto struct {
	BaseResponseDto
	Location string
}

func NewGetResponseDto(statusCode int, contentType string, location string, token string) *GetResponseDto {
	return &GetResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode: statusCode,
		Token:      token,
	}, Location: location}
}
