package middleware

import (
	"github.com/anoriar/shortener/internal/shortener/shared/compress"
	"net/http"
	"strings"
)

const (
	applicationJson = "application/json"
	textHtml        = "text/html"
)

type CompressMiddleware struct {
}

func NewCompressMiddleware() *CompressMiddleware {
	return &CompressMiddleware{}
}

func (cm *CompressMiddleware) Compress(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		supportCompress := false
		contentType := r.Header.Get("Content-Type")
		if contentType == applicationJson || contentType == textHtml {
			supportCompress = true
		}

		if supportCompress == true {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := compress.NewCompressWriter(w)
				ow = cw
				defer cw.Close()
			}

		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := compress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		h.ServeHTTP(ow, r)
	})
}
