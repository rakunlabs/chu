package httploader

import (
	"io"
	"strings"

	"github.com/rakunlabs/chu/utils/decoder"
)

type Decoder = func(r io.Reader, to any) error

func GetDecoder(contentType string) Decoder {
	contentType = strings.Split(strings.ToLower(contentType), ";")[0]

	switch contentType {
	case "application/json":
		return decoder.DecodeJson
	case "application/x-toml":
		return decoder.DecodeToml
	case "application/x-yaml", "text/yaml", "text/x-yaml":
		return decoder.DecodeYaml
	default:
		return decoder.DecodeYaml
	}
}
