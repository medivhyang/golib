package config

type Source interface {
	Encoder
	Decoder
}

type Encoder interface {
	Encode(i interface{}) ([]byte, error)
}

type Decoder interface {
	Decode(source []byte, target interface{}) error
}

type EncoderFunc func(i interface{}) ([]byte, error)

func (f EncoderFunc) Encode(i interface{}) ([]byte, error) {
	if f == nil {
		return nil, ErrNilEncoder
	}
	return f(i)
}

type DecoderFunc func(source []byte, target interface{}) error

func (f DecoderFunc) Decode(source []byte, target interface{}) error {
	if f == nil {
		return ErrNilDecoder
	}
	return f(source, target)
}
