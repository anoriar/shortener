package response

type BaseResponseDto struct {
	StatusCode  int
	ContentType string
	Token       string
}
