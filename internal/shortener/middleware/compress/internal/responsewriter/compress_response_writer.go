package responsewriter

import (
	"compress/gzip"
	"net/http"
)

type compressResponseWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressResponseWriter {
	return &compressResponseWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressResponseWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressResponseWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressResponseWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close закрывает compress.Writer и досылает все данные из буфера.
func (c *compressResponseWriter) Close() error {
	return c.zw.Close()
}
