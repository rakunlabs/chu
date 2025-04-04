package fileloader

import (
	"errors"
	"fmt"
	"io"

	"github.com/rakunlabs/chu/utils/decoderfile"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

type Decoder interface {
	Decode(r io.Reader, to any) error
}

func Decoders() map[string]Decoder {
	yamlDecoder := &decoderfile.Yaml{}

	return map[string]Decoder{
		".toml": &decoderfile.Toml{},
		".yaml": yamlDecoder,
		".yml":  yamlDecoder,
		".json": &decoderfile.Json{},
	}
}

func (l Loader) getFileDecoder(ext string) (Decoder, error) {
	if decoder, ok := l.Decoders[ext]; ok {
		return decoder, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileFormat, ext)
}
