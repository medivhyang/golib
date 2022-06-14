package config

import (
	"encoding/xml"
)

type XMLSource struct {
	XMLEncoder
	XMLDecoder
}

type XMLEncoder struct {
	Pretty bool
	Prefix string
	Indent string
}

func (e XMLEncoder) Encode(i interface{}) ([]byte, error) {
	if e.Pretty {
		return xml.MarshalIndent(i, e.Prefix, e.Indent)
	}
	return xml.Marshal(i)
}

type XMLDecoder struct{}

func (d XMLDecoder) Decode(source []byte, target interface{}) error {
	return xml.Unmarshal(source, target)
}
