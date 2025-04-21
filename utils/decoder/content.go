package decoder

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

func DecodeJson(r io.Reader, to any) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}

func DecodeToml(r io.Reader, to any) error {
	decoder := toml.NewDecoder(r)

	if _, err := decoder.Decode(to); err != nil {
		return fmt.Errorf("toml decoder: %w", err)
	}

	return nil
}

func DecodeYaml(r io.Reader, to any) error {
	decoder := yaml.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("yaml decoder: %w", err)
	}

	return nil
}
