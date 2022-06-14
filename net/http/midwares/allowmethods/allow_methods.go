package allowmethods

import "github.com/medivhyang/golib/net/http"

func New(methods ...string) http.Midware {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w *http.ResponseWriter, r *http.Request) {
			find := false
			for _, method := range methods {
				if r.Method() == method {
					find = true
				}
			}
			if !find {
				w.StatusCode(http.StatusMethodNotAllowed)
				return
			}
			h.Next(w, r)
		}
	}
}
