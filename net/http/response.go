package http

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
)

type Response struct {
	raw    *http.Response
	read   bool
	errors []error
}

func NewResponse(r *http.Response) *Response {
	return &Response{raw: r}
}

func newErrorResponse(err error) *Response {
	if err == nil {
		return &Response{errors: []error{ErrUnknownError}}
	}
	return &Response{errors: []error{err}}
}

func (r *Response) StatusCode() int {
	return r.raw.StatusCode
}

func (r *Response) Dump(body bool) ([]byte, error) {
	return httputil.DumpResponse(r.raw, body)
}

func (r *Response) Pipe(writer io.Writer) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	if _, err := io.Copy(writer, r.raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *Response) SaveFile(filename string) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	if _, err := io.Copy(file, r.raw.Body); err != nil {
		return err
	}
	return nil
}

func (r *Response) Stream() (io.ReadCloser, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	return r.raw.Body, nil
}

func (r *Response) Bytes() ([]byte, error) {
	if r.read {
		return nil, ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	content, err := ioutil.ReadAll(r.raw.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (r *Response) Text() (string, error) {
	if r.read {
		return "", ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	bs, err := r.Bytes()
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (r *Response) JSON(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return json.NewDecoder(r.raw.Body).Decode(value)
}

func (r *Response) XML(value interface{}) error {
	if r.read {
		return ErrResponseBodyHasRead
	}
	defer func() {
		r.raw.Body.Close()
		r.read = true
	}()
	return xml.NewDecoder(r.raw.Body).Decode(value)
}
