package responsewriter

import "net/http"

type ResponseData struct {
	status int
	size   int
}

func (r ResponseData) Status() int {
	return r.status
}

func (r ResponseData) Size() int {
	return r.size
}

func NewResponseData() *ResponseData {
	return &ResponseData{
		0, 0,
	}
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	responseData *ResponseData
}

func (l *LoggingResponseWriter) ResponseData() *ResponseData {
	return l.responseData
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		w,
		NewResponseData(),
	}
}

func (l *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	size, err := l.ResponseWriter.Write(bytes)
	l.responseData.size = size
	return size, err
}

func (l *LoggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}
