package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GZip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")

			gzw := gzip.NewWriter(w)
			defer gzw.Close()

			gzipResponseWriter := &gzipResponseWriter{Writer: gzw, ResponseWriter: w}
			next.ServeHTTP(gzipResponseWriter, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(p []byte) (n int, err error) {
	return g.Writer.Write(p)
}
