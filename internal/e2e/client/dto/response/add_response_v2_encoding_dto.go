package response

type AddResponseV2EncodingDto struct {
	AddResponseV2Dto
	ContentEncoding string
}

func NewAddResponseV2EncodingDto(statusCode int, contentType string, contentEncoding string, body AddURLResponseDTO) *AddResponseV2EncodingDto {
	return &AddResponseV2EncodingDto{
		AddResponseV2Dto: AddResponseV2Dto{
			StatusCode:  statusCode,
			ContentType: contentType,
			Body:        body,
		},
		ContentEncoding: contentEncoding,
	}
}
