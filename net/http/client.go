package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	ErrOpenFileFailed        = errorf("open file failed")
	ErrTemplateRequireURL    = errorf("template require url")
	ErrTemplateRequireMethod = errorf("template require method")
)

func Get(path string) *Response {
	return NewBuilder().Get(path).Do()
}

func GetText(path string) (string, error) {
	return NewBuilder().Get(path).Do().Text()
}

func GetJSON(path string, result interface{}) error {
	return NewBuilder().Get(path).Do().JSON(result)
}

func SaveFile(path string, filename string) error {
	return NewBuilder().Get(path).Do().SaveFile(filename)
}

func GetStream(path string) (io.ReadCloser, error) {
	return NewBuilder().Get(path).Do().Stream()
}

func Post(path string, reader io.Reader, contentType string) *Response {
	return NewBuilder().Post(path).Header("Content-Type", contentType).Write(reader).Do()
}

func PostJSON(path string, body interface{}, result interface{}) error {
	return NewBuilder().Post(path).WriteJSON(body).Do().JSON(result)
}

func PostFile(path string, fileName string) *Response {
	return NewBuilder().Post(path).WriteFile(fileName).Do()
}

func PostFormFile(path string, formName string, fileName string) *Response {
	return NewBuilder().Post(path).WriteFormFile(formName, fileName).Do()
}

func Put(path string, reader io.Reader, contentType string) *Response {
	return NewBuilder().Put(path).Header("Content-Type", contentType).Write(reader).Do()
}

func PutJSON(path string, body interface{}, result interface{}) error {
	return NewBuilder().Post(path).WriteJSON(body).Do().JSON(result)
}

func Patch(path string, reader io.Reader, contentType string) *Response {
	return NewBuilder().Patch(path).Header("Content-Type", contentType).Write(reader).Do()
}

func PatchJSON(path string, body interface{}, result interface{}) error {
	return NewBuilder().Post(path).WriteJSON(body).Do().JSON(result)
}

type Template struct {
	Prefix  string
	Method  string
	Path    string
	Queries map[string]string
	Headers map[string]string
	Body    io.Reader
}

func (t *Template) New() *Builder {
	return NewBuilder(t)
}

func (t *Template) Do(client ...*http.Client) *Response {
	if t.Method == "" {
		return newErrorResponse(ErrTemplateRequireMethod)
	}

	fullURL := t.FullURL()
	if fullURL == "" {
		return newErrorResponse(ErrTemplateRequireURL)
	}

	request, err := http.NewRequest(t.Method, fullURL, t.Body)
	if err != nil {
		return newErrorResponse(err)
	}

	for k, v := range t.Headers {
		request.Header.Set(k, v)
	}

	var finalClient *http.Client
	if len(client) > 0 {
		finalClient = client[0]
	}
	if finalClient == nil {
		finalClient = http.DefaultClient
	}

	response, err := finalClient.Do(request)
	if err != nil {
		return newErrorResponse(err)
	}

	return NewResponse(response)
}

func (t *Template) FullURL() string {
	var result string
	if t.Prefix != "" && !strings.HasSuffix(t.Prefix, "/") && t.Path != "" && !strings.HasPrefix(t.Path, "/") {
		result = t.Prefix + "/" + t.Path
	} else if strings.HasSuffix(t.Prefix, "/") && strings.HasPrefix(t.Path, "/") {
		result = strings.TrimSuffix(t.Prefix, "/") + t.Path
	} else {
		result = t.Prefix + t.Path
	}
	if strings.Contains(result, "?") {
		result += "&"
	} else {
		result += "?"
	}
	if len(t.Queries) > 0 {
		vs := url.Values{}
		for k, v := range t.Queries {
			vs.Set(k, v)
		}
		result += vs.Encode()
	}
	return result
}

func (t *Template) copy() (*Template, error) {
	result := Template{
		Prefix:  t.Prefix,
		Method:  t.Method,
		Path:    t.Path,
		Queries: make(map[string]string, len(t.Queries)),
		Headers: make(map[string]string, len(t.Headers)),
	}
	for k, v := range t.Queries {
		result.Queries[k] = v
	}
	for k, v := range t.Headers {
		result.Queries[k] = v
	}
	if t.Body != nil {
		buf := &bytes.Buffer{}
		if _, err := io.Copy(buf, t.Body); err != nil {
			return nil, err
		}
	}
	return &result, nil
}

type Builder struct {
	Template *Template
	errors   []error
}

func NewBuilder(t ...*Template) *Builder {
	var finalT *Template
	if len(t) > 0 && t[0] != nil {
		c, err := t[0].copy()
		if err != nil {
			return &Builder{errors: []error{err}}
		}
		finalT = c
	}
	if finalT == nil {
		finalT = &Template{}
	}
	return &Builder{Template: finalT}
}

func (b *Builder) Prefix(p string) *Builder {
	b.Template.Prefix = p
	return b
}

func (b *Builder) Get(path string) *Builder {
	b.Template.Method = http.MethodGet
	b.Template.Path = path
	return b
}

func (b *Builder) Post(path string) *Builder {
	b.Template.Method = http.MethodPost
	b.Template.Path = path
	return b
}

func (b *Builder) Put(path string) *Builder {
	b.Template.Method = http.MethodPut
	b.Template.Path = path
	return b
}

func (b *Builder) Delete(path string) *Builder {
	b.Template.Method = http.MethodDelete
	b.Template.Path = path
	return b
}

func (b *Builder) Patch(path string) *Builder {
	b.Template.Method = http.MethodPatch
	b.Template.Path = path
	return b
}

func (b *Builder) Options(path string) *Builder {
	b.Template.Method = http.MethodOptions
	b.Template.Path = path
	return b
}

func (b *Builder) Query(k, v string) *Builder {
	b.Template.Queries[k] = v
	return b
}

func (b *Builder) Queries(m map[string]string) *Builder {
	if b.Template.Queries == nil {
		b.Template.Queries = map[string]string{}
	}
	for k, v := range m {
		b.Template.Queries[k] = v
	}
	return b
}

func (b *Builder) CleanQueries() *Builder {
	b.Template.Queries = map[string]string{}
	return b
}

func (b *Builder) Header(k, v string) *Builder {
	if b.Template.Headers == nil {
		b.Template.Headers = map[string]string{}
	}
	b.Template.Headers[k] = v
	return b
}

func (b *Builder) Headers(m map[string]string) *Builder {
	for k, v := range m {
		b.Header(k, v)
	}
	return b
}

func (b *Builder) CleanHeaders() *Builder {
	b.Template.Headers = map[string]string{}
	return b
}

func (b *Builder) Write(reader io.Reader) *Builder {
	b.Template.Body = reader
	return b
}

func (b *Builder) WriteText(text string) *Builder {
	b.Template.Body = strings.NewReader(text)
	return b
}

func (b *Builder) WriteJSON(v interface{}) *Builder {
	bs, err := json.Marshal(v)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	b.Template.Body = bytes.NewReader(bs)
	return b
}

func (b *Builder) WriteXML(v interface{}) *Builder {
	bs, err := xml.Marshal(v)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	b.Template.Body = bytes.NewReader(bs)
	return b
}

func (b *Builder) WriteFile(filename string) *Builder {
	file, err := os.Open(filename)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	b.Template.Body = file
	return b
}

func (b *Builder) WriteFormFile(formName string, fileName string) *Builder {
	w := multipart.NewWriter(&bytes.Buffer{})
	defer w.Close()

	ff, err := w.CreateFormFile(formName, fileName)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	file, err := os.Open(fileName)
	if err != nil {
		b.errors = append(b.errors, ErrOpenFileFailed)
		return b
	}
	defer file.Close()
	if _, err = io.Copy(ff, file); err != nil {
		b.errors = append(b.errors, ErrOpenFileFailed)
		return b
	}

	b.Header("Content-Type", w.FormDataContentType())

	return b
}

func (b *Builder) CleanBody() *Builder {
	b.Template.Body = nil
	return b
}

func (b *Builder) Build() *Template {
	return b.Template
}

func (b *Builder) Do(client ...*http.Client) *Response {
	if len(b.errors) > 0 {
		return &Response{errors: b.errors}
	}
	return b.Template.Do(client...)
}

// endregion
