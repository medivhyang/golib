package config

import (
	"bytes"

	"github.com/BurntSushi/toml"
)

type TOMLSource struct {
	TOMLEncoder
	TOMLDecoder
}

type TOMLEncoder struct {
	Pretty bool
	Indent string
}

func (e TOMLEncoder) Encode(i interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	encoder := toml.NewEncoder(&buf)
	if e.Pretty {
		encoder.Indent = e.Indent
	}
	if err := encoder.Encode(i); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type TOMLDecoder struct{}

func (d TOMLDecoder) Decode(source []byte, i interface{}) error {
	return toml.Unmarshal(source, i)
}
