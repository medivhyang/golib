package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/medivhyang/golib/net/http"
)

var Default = New()

type HandleFunc func(w *http.ResponseWriter, r *http.Request, err interface{})

func New(callback ...HandleFunc) http.Midware {
	var f func(w *http.ResponseWriter, r *http.Request, err interface{})
	if len(callback) > 0 {
		f = callback[0]
	} else {
		f = DefaultHandleFunc
	}
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w *http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					f(w, r, err)
				}
			}()
			h.ServeHTTP(w.Raw(), r.Raw())
		}
	}
}

func DefaultHandleFunc(w *http.ResponseWriter, r *http.Request, err interface{}) {
	fmt.Printf("recovery: %+v\n%s", err, string(debug.Stack()))
	w.Text(http.StatusInternalServerError, fmt.Sprint(err))
}
