package trace

import (
	"log"

	"github.com/medivhyang/golib/net/http"
)

var Default = New()

type HandleFunc func(w *http.ResponseWriter, r *http.Request)

func New(callback ...HandleFunc) http.Midware {
	var f HandleFunc
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = DefaultHandleFunc
	}
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w *http.ResponseWriter, r *http.Request) {
			f(w, r)
			h.Next(w, r)
		}
	}
}

func DefaultHandleFunc(w *http.ResponseWriter, r *http.Request) {
	log.Printf("trace: %s %s", r.Method(), r.Path())
}
