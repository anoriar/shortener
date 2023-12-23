package response

// AddURLResponseDTO missing godoc.
type AddURLResponseDTO struct {
	Result string `json:"result"`
}

// AddResponseV2Dto missing godoc.
type AddResponseV2Dto struct {
	BaseResponseDto
	Body AddURLResponseDTO
}

// NewAddResponseV2Dto missing godoc.
func NewAddResponseV2Dto(statusCode int, contentType string, token string, body AddURLResponseDTO) *AddResponseV2Dto {
	return &AddResponseV2Dto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Body: body}
}
