package response

// PingResponseDto missing godoc.
type PingResponseDto struct {
	StatusCode int
}

// NewPingResponseDto missing godoc.
func NewPingResponseDto(statusCode int) PingResponseDto {
	return PingResponseDto{StatusCode: statusCode}
}
