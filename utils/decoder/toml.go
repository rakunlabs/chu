package decoder

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type Toml struct{}

func (Toml) Decode(r io.Reader, to interface{}) error {
	decoder := toml.NewDecoder(r)

	if _, err := decoder.Decode(to); err != nil {
		return fmt.Errorf("toml decoder: %w", err)
	}

	return nil
}
