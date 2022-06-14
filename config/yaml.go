package config

import (
	"gopkg.in/yaml.v2"
)

type YAMLSource struct {
	YAMLEncoder
	YAMLDecoder
}

type YAMLEncoder struct{}

func (e YAMLEncoder) Encode(i interface{}) ([]byte, error) {
	return yaml.Marshal(i)
}

type YAMLDecoder struct{}

func (d YAMLDecoder) Decode(source []byte, target interface{}) error {
	return yaml.Unmarshal(source, target)
}
