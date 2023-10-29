package response

type GetUserURLsResponseDto struct {
	BaseResponseDto
	Items []UserURLResponseItem
}

type UserURLResponseItem struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewGetUserURLsResponseDto(statusCode int, contentType string, token string, items []UserURLResponseItem) *GetUserURLsResponseDto {
	return &GetUserURLsResponseDto{BaseResponseDto: BaseResponseDto{
		StatusCode:  statusCode,
		ContentType: contentType,
		Token:       token,
	}, Items: items}
}
