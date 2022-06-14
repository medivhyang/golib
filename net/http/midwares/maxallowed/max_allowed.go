package maxallowed

import "github.com/medivhyang/golib/net/http"

func New(n int) http.Midware {
	sem := make(chan struct{}, n)
	acquire := func() { sem <- struct{}{} }
	release := func() { <-sem }
	return func(f http.HandlerFunc) http.HandlerFunc {
		acquire()
		defer release()
		return func(w *http.ResponseWriter, r *http.Request) {
			f.Next(w, r)
		}
	}
}
