package response

type PingResponseDto struct {
	StatusCode int
}

func NewPingResponseDto(statusCode int) *PingResponseDto {
	return &PingResponseDto{StatusCode: statusCode}
}
