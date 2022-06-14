package config

import "encoding/json"

type JSONSource struct {
	JSONEncoder
	JSONDecoder
}

type JSONEncoder struct {
	Pretty bool
	Prefix string
	Indent string
}

func (e JSONEncoder) Encode(i interface{}) ([]byte, error) {
	if e.Pretty {
		return json.MarshalIndent(i, e.Prefix, e.Indent)
	}
	return json.Marshal(i)
}

type JSONDecoder struct{}

func (d JSONDecoder) Decode(source []byte, target interface{}) error {
	return json.Unmarshal(source, target)
}
