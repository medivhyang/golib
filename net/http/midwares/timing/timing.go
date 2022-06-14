package timing

import (
	"log"
	"time"

	"github.com/medivhyang/golib/net/http"
)

var Default = New()

type HandleFunc func(w *http.ResponseWriter, r *http.Request, d time.Duration)

func New(callback ...HandleFunc) http.Midware {
	var f HandleFunc
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = DefaultHandleFunc
	}
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w *http.ResponseWriter, r *http.Request) {
			defer func(start time.Time) {
				d := time.Since(start)
				f(w, r, d)
			}(time.Now())
			h.Next(w, r)
		}
	}
}

func DefaultHandleFunc(w *http.ResponseWriter, r *http.Request, d time.Duration) {
	log.Printf("timing: \"%s %s\" cost %s\n", r.Method(), r.Path(), d)
}
