package basicauth

import (
	"errors"
	"fmt"

	"github.com/medivhyang/golib/net/http"
)

var ErrRequireFunc = errors.New("require func")

func New(pairs map[string]string, realm ...string) http.Midware {
	return NewWithFunc(func(username, password string) bool {
		for k, v := range pairs {
			if username == k && password == v {
				return true
			}
		}
		return false
	}, realm...)
}

func NewWithFunc(f func(username, password string) bool, realm ...string) http.Midware {
	if f == nil {
		panic(ErrRequireFunc)
	}
	var finalRealm string
	if len(realm) > 0 {
		finalRealm = realm[0]
	} else {
		finalRealm = ""
	}
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w *http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if ok {
				if f(username, password) {
					h.ServeHTTP(w.Raw(), r.Raw())
					return
				}
			}
			w.Header("WWW-Authenticate", fmt.Sprintf("Basic realm=%q", finalRealm))
			w.StatusCode(http.StatusUnauthorized)
			return
		}
	}
}
