package fileloader

import (
	"errors"
	"fmt"
	"io"

	"github.com/rakunlabs/chu/utils/decoder"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

type Decoder = func(r io.Reader, to any) error

func Decoders() map[string]Decoder {
	return map[string]Decoder{
		".toml": decoder.DecodeToml,
		".yaml": decoder.DecodeYaml,
		".yml":  decoder.DecodeYaml,
		".json": decoder.DecodeJson,
	}
}

func (l Loader) getFileDecoder(ext string) (Decoder, error) {
	if decoder, ok := l.Decoders[ext]; ok {
		return decoder, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileFormat, ext)
}
