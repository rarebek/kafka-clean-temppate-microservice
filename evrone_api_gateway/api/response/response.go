package response

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
	code int
}

func NewResponseWriter(w http.ResponseWriter, statusCode int) *ResponseWriter {
	return &ResponseWriter{w, statusCode}
}

func (w *ResponseWriter) StatusCode() int {
	return w.code
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("given http.ResponseWriter is not a http.Hijacker")
	}
	return h.Hijack()
}
