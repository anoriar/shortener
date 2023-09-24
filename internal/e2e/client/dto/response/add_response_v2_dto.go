package response

type AddURLResponseDTO struct {
	Result string `json:"result"`
}

type AddResponseV2Dto struct {
	StatusCode  int
	ContentType string
	Body        AddURLResponseDTO
}

func NewAddResponseV2Dto(statusCode int, contentType string, body AddURLResponseDTO) *AddResponseV2Dto {
	return &AddResponseV2Dto{StatusCode: statusCode, ContentType: contentType, Body: body}
}
