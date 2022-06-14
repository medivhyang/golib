package http

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/medivhyang/golib/reflect/binding"
)

type Request struct {
	raw    *http.Request
	params map[string]string
}

func NewRequest(r *http.Request) *Request {
	return &Request{raw: r}
}

func (r *Request) Raw() *http.Request {
	return r.raw
}

func (r *Request) Context() context.Context {
	return r.raw.Context()
}

func (r *Request) SetContext(ctx context.Context) {
	r.raw = r.raw.WithContext(ctx)
}

func (r *Request) Method() string {
	return r.raw.Method
}

func (r *Request) Path() string {
	return r.raw.URL.Path
}

func (r *Request) Header(key string) string {
	return r.raw.Header.Get(key)
}

func (r *Request) HeaderOrDefault(key string, value string) string {
	result := r.raw.Header.Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) HeaderExists(key string) bool {
	if r.raw.Header == nil {
		return false
	}
	return len(r.raw.Header[key]) > 0
}

func (r *Request) SetHeader(key string, value string) {
	r.raw.Header.Set(key, value)
}

func (r *Request) Cookie(key string) (string, error) {
	cookie, err := r.raw.Cookie(key)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (r *Request) AddCooke(cookie *http.Cookie) {
	r.raw.AddCookie(cookie)
}

func (r *Request) CookieOrDefault(key string, defaultValue ...string) string {
	cookie, err := r.raw.Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			if len(defaultValue) > 0 {
				return defaultValue[0]
			}
			return ""
		}
		panic(err)
	}
	return cookie.Value
}

func (r *Request) CookieExists(key string) bool {
	_, err := r.raw.Cookie(key)
	if err != nil {
		if err == http.ErrNoCookie {
			return false
		}
		panic(err)
	}
	return true
}

func (r *Request) Param(key string) string {
	if r.params == nil {
		r.params = Params(r.raw)
	}
	return r.params[key]
}

func (r *Request) ParamOrDefault(key string, value string) string {
	if r.params == nil {
		r.params = Params(r.raw)
	}
	result := r.params[key]
	if result == "" {
		return value
	}
	return result
}

func (r *Request) ParamExists(key string) bool {
	if r.params == nil {
		r.params = Params(r.raw)
	}
	_, ok := r.params[key]
	return ok
}

func (r *Request) Query(key string) string {
	return r.raw.URL.Query().Get(key)
}

func (r *Request) QueryOrDefault(key string, value string) string {
	result := r.raw.URL.Query().Get(key)
	if result == "" {
		return value
	}
	return result
}

func (r *Request) QueryExists(key string) bool {
	values := r.raw.URL.Query()
	if values == nil {
		return false
	}
	return len(values[key]) > 0
}

func (r *Request) PostForm(key string) string {
	return r.raw.PostFormValue(key)
}

func (r *Request) PostFormExists(key string) bool {
	_ = r.raw.ParseForm()
	if r.raw.PostForm == nil {
		return false
	}
	return len(r.raw.PostForm[key]) > 0
}

func (r *Request) Form(key string) string {
	return r.raw.FormValue(key)
}

func (r *Request) FormExists(key string) bool {
	_ = r.raw.ParseForm()
	if r.raw.Form == nil {
		return false
	}
	return len(r.raw.Form[key]) > 0
}

func (r *Request) BindJSON(i interface{}) error {
	return json.NewDecoder(r.raw.Body).Decode(i)
}

func (r *Request) BindXML(i interface{}) error {
	return xml.NewDecoder(r.raw.Body).Decode(i)
}

func (r *Request) BindQuery(i interface{}) error {
	m := parseStructTagFirstWord(i, BindingTagKey)
	return binding.BindStructFunc(func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = toCase(caseSnake, s)
		}
		return r.raw.URL.Query()[s]
	}, i)
}

func (r *Request) BindPostForm(i interface{}) error {
	if err := r.raw.ParseForm(); err != nil {
		return err
	}
	m := parseStructTagFirstWord(i, BindingTagKey)
	return binding.BindStructFunc(func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = toCase(caseSnake, s)
		}
		return r.raw.PostForm[s]
	}, i)
}

func (r *Request) BindForm(i interface{}) error {
	if err := r.raw.ParseForm(); err != nil {
		return err
	}
	m := parseStructTagFirstWord(i, BindingTagKey)
	return binding.BindStructFunc(func(s string) []string {
		if v, ok := m[s]; ok {
			s = v
		} else {
			s = toCase(caseSnake, s)
		}
		return r.raw.Form[s]
	}, i)
}

func (r *Request) BasicAuth() (username string, password string, ok bool) {
	return r.raw.BasicAuth()
}

func (r *Request) SetBasicAuth(username string, password string) {
	r.raw.SetBasicAuth(username, password)
}
