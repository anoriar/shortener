package response

// AddResponseV2EncodingDto missing godoc.
type AddResponseV2EncodingDto struct {
	AddResponseV2Dto
	ContentEncoding string
}

// NewAddResponseV2EncodingDto missing godoc.
func NewAddResponseV2EncodingDto(statusCode int, contentType string, token string, contentEncoding string, body AddURLResponseDTO) *AddResponseV2EncodingDto {
	return &AddResponseV2EncodingDto{
		AddResponseV2Dto: AddResponseV2Dto{
			BaseResponseDto: BaseResponseDto{
				StatusCode:  statusCode,
				ContentType: contentType,
				Token:       token,
			},
			Body: body,
		},

		ContentEncoding: contentEncoding,
	}
}
