package http

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
)

type Context interface {
	Request() *http.Request
	SetRequest(r *http.Request)
	RresponseWriter() http.ResponseWriter
	SetResponseWriter(w http.ResponseWriter)
	Scheme() string
	Method() string
	Path() string
	Param(string) string
	Query(string) string
	Header(string) string
	PostForm(string) string
	Cookie(string) (*http.Cookie, error)
	Bind(i interface{}) error
	IsTLS() bool
	RealIP() string

	Text(code int, text string) error
	HTML(code int, html string) error
	JSON(code int, i interface{}) error
}

type contextImpl struct {
	request        *http.Request
	responseWriter http.ResponseWriter
}

func (c *contextImpl) Request() *http.Request {
	return c.request
}

func (c *contextImpl) SetRequest(r *http.Request) {
	c.request = r
}

func (c *contextImpl) ResponseWriter() http.ResponseWriter {
	return c.responseWriter
}

func (c *contextImpl) SetResponseWriter(w http.ResponseWriter) {
	c.responseWriter = w
}

func (c *contextImpl) IsTLS() bool {
	return c.request.TLS != nil
}

func (c *contextImpl) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.request.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}
func (c *contextImpl) Method() string { return c.request.Method }

func (c *contextImpl) Path() string {
	return c.request.URL.Path
}

func (c *contextImpl) Param(name string) string {
	panic("not implemented")
}

func (c *contextImpl) Query(name string) string {
	return c.request.URL.Query().Get(name)
}

func (c *contextImpl) Header(name string) string {
	return c.request.Header.Get(name)
}

func (c *contextImpl) PostForm(name string) string {
	return c.request.PostFormValue(name)
}

func (c *contextImpl) FormFile(name string) (*multipart.FileHeader, error) {
	f, fh, err := c.request.FormFile(name)
	if err != nil {
		return nil, err
	}
	f.Close()
	return fh, nil
}

func (c *contextImpl) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

func (c *contextImpl) Bind(i interface{}) error {
	panic("not implemented")
}

func (c *contextImpl) Text(code int, text string) error {
	w := c.ResponseWriter()
	w.Header().Set(HeaderContentType, MIMEText)
	w.WriteHeader(code)
	if _, err := io.WriteString(w, text); err != nil {
		return err
	}
	return nil
}

func (c *contextImpl) HTML(code int, html string) error {
	w := c.ResponseWriter()
	w.Header().Set(HeaderContentType, MIMEHTML)
	w.WriteHeader(code)
	if _, err := io.WriteString(w, html); err != nil {
		return err
	}
	return nil
}

func (c *contextImpl) JSON(code int, i interface{}) error {
	w := c.ResponseWriter()
	w.Header().Set(HeaderContentType, MIMEJSON)
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(i); err != nil {
		return err
	}
	return nil
}

func (c *contextImpl) RealIP() string {
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			return strings.TrimSpace(ip[:i])
		}
		return ip
	}
	if ip := c.request.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.request.RemoteAddr)
	return ra
}
