package responsewriter

import (
	"net/http"

	"github.com/klauspost/compress/gzip"
)

type KlauspostGzipCompressResponseWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewKlauspostGzipCompressWriter(w http.ResponseWriter) (*KlauspostGzipCompressResponseWriter, error) {
	zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	return &KlauspostGzipCompressResponseWriter{
		w:  w,
		zw: zw,
	}, nil
}

func (c *KlauspostGzipCompressResponseWriter) Header() http.Header {
	return c.w.Header()
}

func (c *KlauspostGzipCompressResponseWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *KlauspostGzipCompressResponseWriter) WriteHeader(statusCode int) {
	c.w.Header().Set("Content-Encoding", "gzip")
	c.w.WriteHeader(statusCode)
}

// Close закрывает compress.Writer и досылает все данные из буфера.
func (c *KlauspostGzipCompressResponseWriter) Close() error {
	return c.zw.Close()
}
