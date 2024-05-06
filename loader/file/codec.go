package file

import (
	"errors"
	"fmt"
	"io"

	"github.com/rakunlabs/chu/utils/decoder"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

type Decoder interface {
	Decode(r io.Reader, to interface{}) error
}

func getDecoders() map[string]Decoder {
	yamlDecoder := &decoder.Yaml{}

	return map[string]Decoder{
		".toml": &decoder.Toml{},
		".yaml": yamlDecoder,
		".yml":  yamlDecoder,
		".json": &decoder.Json{},
	}
}

func (l Loader) getFileDecoder(ext string) (Decoder, error) {
	if decoder, ok := l.decoders[ext]; ok {
		return decoder, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileFormat, ext)
}
