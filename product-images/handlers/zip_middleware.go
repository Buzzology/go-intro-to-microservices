package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)

	return &WrappedResponseWriter{rw: rw, gw: gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

func (wr *WrappedResponseWriter) Write(data []byte) (int, error) {
	return wr.gw.Write(data)
}

func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
	wr.Flush()
}

// GzipMiddleware To use as middleware we need to convert the handler to a handler func
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {

			// Create a new instance of our wrapped response writer
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			// Send back our wrapped (gzip) response writer
			next.ServeHTTP(wrw, req)

			// Ensure that flush is called
			defer wrw.Flush()
			return
		}

		// No compression
		next.ServeHTTP(rw, req)
	})
}
