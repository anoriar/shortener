package reader

import (
	"io"

	"github.com/klauspost/compress/gzip"
)

// KlauspostGzipCompressReader missing godoc.
type KlauspostGzipCompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

// NewKlauspostGzipCompressReader missing godoc.
func NewKlauspostGzipCompressReader(r io.ReadCloser) (*KlauspostGzipCompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &KlauspostGzipCompressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read missing godoc.
func (c KlauspostGzipCompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close missing godoc.
func (c *KlauspostGzipCompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
