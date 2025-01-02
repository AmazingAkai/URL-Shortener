package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AmazingAkai/URL-Shortener/app/internal/log"
	"github.com/fatih/color"
)

var (
	green   = color.New(color.FgGreen).Sprint
	cyan    = color.New(color.FgCyan).Sprint
	yellow  = color.New(color.FgYellow).Sprint
	blue    = color.New(color.FgBlue).Sprint
	magenta = color.New(color.FgMagenta).Sprint
	red     = color.New(color.FgRed).Sprint
	reset   = color.New(color.Reset).Sprint
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		srw := &statusResponseWriter{ResponseWriter: w}

		start := time.Now()
		next.ServeHTTP(srw, r)
		timeTaken := fmt.Sprintf("%.3f µs", float64(time.Since(start))/1e3)

		log.Infof("%s - %s - %s - %s - %s - %s - %s %s",
			green(r.Method),
			cyan(r.URL.Path),
			blue(r.Proto),
			yellow(r.RemoteAddr),
			red(srw.statusCode),
			magenta(r.ContentLength),
			cyan(timeTaken),
			reset())
	})
}
