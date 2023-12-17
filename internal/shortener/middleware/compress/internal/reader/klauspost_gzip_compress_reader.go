package reader

import (
	"github.com/klauspost/compress/gzip"
	"io"
)

type KlauspostGzipCompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

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

func (c KlauspostGzipCompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *KlauspostGzipCompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
