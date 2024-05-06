package decoder

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Yaml struct{}

func (Yaml) Decode(r io.Reader, to interface{}) error {
	decoder := yaml.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("yaml decoder: %w", err)
	}

	return nil
}
