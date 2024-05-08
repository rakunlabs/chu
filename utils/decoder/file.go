package decoder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Json struct{}

func (Json) Decode(r io.Reader, to interface{}) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}

type Toml struct{}

func (Toml) Decode(r io.Reader, to interface{}) error {
	decoder := toml.NewDecoder(r)

	if _, err := decoder.Decode(to); err != nil {
		return fmt.Errorf("toml decoder: %w", err)
	}

	return nil
}

type Yaml struct{}

func (Yaml) Decode(r io.Reader, to interface{}) error {
	decoder := yaml.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("yaml decoder: %w", err)
	}

	return nil
}
