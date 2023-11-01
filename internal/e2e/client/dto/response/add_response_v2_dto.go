package response

type AddURLResponseDTO struct {
	Result string `json:"result"`
}

type AddResponseV2Dto struct {
	BaseResponseDto
	Body AddURLResponseDTO
}

func NewAddResponseV2Dto(statusCode int, contentType string, token string, body AddURLResponseDTO) *AddResponseV2Dto {
	return &AddResponseV2Dto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
