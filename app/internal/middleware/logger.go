package middleware

import (
	"net/http"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen).Sprint
	cyan   = color.New(color.FgCyan).Sprint
	yellow = color.New(color.FgYellow).Sprint
	red    = color.New(color.FgRed).Sprint
	reset  = color.New(color.Reset).Sprint
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusResponseWriter{ResponseWriter: w}
		next.ServeHTTP(srw, r)

		log.Infof("%s - %s - %s - %s %s",
			green(r.Method),
			cyan(r.URL.Path),
			yellow(r.RemoteAddr),
			red(srw.statusCode),
			reset())
	})
}
