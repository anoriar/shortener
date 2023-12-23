package responsewriter

import "net/http"

// ResponseData missing godoc.
type ResponseData struct {
	status int
	size   int
}

// Status missing godoc.
func (r ResponseData) Status() int {
	return r.status
}

// Size missing godoc.
func (r ResponseData) Size() int {
	return r.size
}

// NewResponseData missing godoc.
func NewResponseData() *ResponseData {
	return &ResponseData{
		0, 0,
	}
}

// LoggingResponseWriter missing godoc.
type LoggingResponseWriter struct {
	http.ResponseWriter
	responseData *ResponseData
}

// ResponseData missing godoc.
func (l *LoggingResponseWriter) ResponseData() *ResponseData {
	return l.responseData
}

// NewLoggingResponseWriter missing godoc.
func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{
		w,
		NewResponseData(),
	}
}

// Write missing godoc.
func (l *LoggingResponseWriter) Write(bytes []byte) (int, error) {
	size, err := l.ResponseWriter.Write(bytes)
	l.responseData.size = size
	return size, err
}

// WriteHeader missing godoc.
func (l *LoggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}
