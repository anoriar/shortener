package response

type GetResponseDto struct {
	StatusCode int
	Location   string
}

func NewGetResponseDto(statusCode int, location string) *GetResponseDto {
	return &GetResponseDto{StatusCode: statusCode, Location: location}
}
