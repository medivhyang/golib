package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
)

type ResponseWriter struct {
	raw http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{raw: w}
}

func (w *ResponseWriter) Raw() http.ResponseWriter {
	return w.raw
}

func (w *ResponseWriter) StatusCode(statusCode int) {
	w.raw.WriteHeader(statusCode)
}

func (w *ResponseWriter) Header(key string, value string) *ResponseWriter {
	w.raw.Header().Set(key, value)
	return w
}

func (w *ResponseWriter) Text(statusCode int, text string) {
	w.raw.Header().Set(HeaderContentType, MIMEText)
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, text); err != nil {
		panic(err)
	}
}

func (w *ResponseWriter) HTML(statusCode int, content string) {
	w.raw.Header().Set(HeaderContentType, MIMEHTML)
	w.raw.WriteHeader(statusCode)
	if _, err := io.WriteString(w.raw, content); err != nil {
		panic(err)
	}
}

func (w *ResponseWriter) JSON(statusCode int, value interface{}) {
	w.raw.Header().Set(HeaderContentType, MIMEJSON)
	w.raw.WriteHeader(statusCode)
	if err := json.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}

func (w *ResponseWriter) XML(statusCode int, value interface{}) {
	w.raw.Header().Set(HeaderContentType, MIMEXML)
	w.raw.WriteHeader(statusCode)
	if err := xml.NewEncoder(w.raw).Encode(value); err != nil {
		panic(err)
	}
}
