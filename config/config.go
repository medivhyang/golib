package config

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrNilDecoder = errors.New("config: nil decoder")
	ErrNilEncoder = errors.New("config: nil encoder")
	ErrNilSource  = errors.New("config: nil source")

	DefaultJSONEncoder = JSONEncoder{}
	DefaultJSONDecoder = JSONDecoder{}
	DefaultJSONSource  = JSONSource{JSONEncoder: DefaultJSONEncoder, JSONDecoder: DefaultJSONDecoder}

	DefaultXMLEncoder = XMLEncoder{}
	DefaultXMLDecoder = XMLDecoder{}
	DefaultXMLSource  = XMLSource{XMLEncoder: DefaultXMLEncoder, XMLDecoder: DefaultXMLDecoder}

	DefaultYAMLEncoder = YAMLEncoder{}
	DefaultYAMLDecoder = YAMLDecoder{}
	DefaultYAMLSource  = YAMLSource{YAMLEncoder: DefaultYAMLEncoder, YAMLDecoder: DefaultYAMLDecoder}

	DefaultTOMLEncoder = TOMLEncoder{}
	DefaultTOMLDecoder = TOMLDecoder{}
	DefaultTOMLSource  = TOMLSource{TOMLEncoder: DefaultTOMLEncoder, TOMLDecoder: DefaultTOMLDecoder}
)

func Load(decoder Decoder, reader io.Reader, target interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	content, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return decoder.Decode(content, target)
}

func Store(encoder Encoder, i interface{}, writer io.Writer) error {
	if encoder == nil {
		return ErrNilEncoder
	}
	content, err := encoder.Encode(i)
	if err != nil {
		return err
	}
	if _, err = writer.Write(content); err != nil {
		return err
	}
	return nil
}

func LoadBytes(decoder Decoder, source []byte, target interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	return decoder.Decode(source, target)
}

func LoadString(decoder Decoder, source string, target interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	return decoder.Decode([]byte(source), target)
}

func LoadEnv(decoder Decoder, key string, target interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	return decoder.Decode([]byte(os.Getenv(key)), target)
}

func StoreEnv(encoder Encoder, i interface{}, key string) error {
	if encoder == nil {
		return ErrNilEncoder
	}
	content, err := encoder.Encode(i)
	if err != nil {
		return err
	}
	return os.Setenv(key, string(content))
}

func LoadOrStoreEnv(source Source, key string, i interface{}) error {
	if source == nil {
		return ErrNilSource
	}
	if _, ok := os.LookupEnv(key); ok {
		return LoadEnv(source, key, i)
	}
	if err := StoreEnv(source, i, key); err != nil {
		return err
	}
	return nil
}

func LoadFile(decoder Decoder, filePath string, target interface{}) error {
	if decoder == nil {
		return ErrNilDecoder
	}
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return os.ErrNotExist
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	return decoder.Decode(content, target)
}

func StoreFile(encoder Encoder, i interface{}, filePath string) error {
	if encoder == nil {
		return ErrNilEncoder
	}
	content, err := encoder.Encode(i)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, content, os.ModePerm)
}

func LoadOrStoreFile(source Source, filePath string, value interface{}) error {
	if source == nil {
		return ErrNilSource
	}
	_, err := os.Stat(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return StoreFile(source, value, filePath)
	}
	return LoadFile(source, filePath, value)
}
