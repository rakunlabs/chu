package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

var ErrUnsupportedFileFormat = errors.New("unsupported file format")

type Decoder interface {
	Decode(r io.Reader, to interface{}) error
}

func getDecoders() map[string]Decoder {
	yamlDecoder := &yamlDecoder{}

	return map[string]Decoder{
		".toml": &tomlDecoder{},
		".yaml": yamlDecoder,
		".yml":  yamlDecoder,
		".json": &jsonDecoder{},
	}
}

func (l Loader) getFileDecoder(ext string) (Decoder, error) {
	if decoder, ok := l.decoders[ext]; ok {
		return decoder, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedFileFormat, ext)
}

type tomlDecoder struct{}

func (tomlDecoder) Decode(r io.Reader, to interface{}) error {
	decoder := toml.NewDecoder(r)

	if _, err := decoder.Decode(to); err != nil {
		return fmt.Errorf("toml decoder: %w", err)
	}

	return nil
}

type yamlDecoder struct{}

func (yamlDecoder) Decode(r io.Reader, to interface{}) error {
	decoder := yaml.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("yaml decoder: %w", err)
	}

	return nil
}

type jsonDecoder struct{}

func (jsonDecoder) Decode(r io.Reader, to interface{}) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}
