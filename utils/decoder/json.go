package decoder

import (
	"encoding/json"
	"fmt"
	"io"
)

type Json struct{}

func (Json) Decode(r io.Reader, to interface{}) error {
	decoder := json.NewDecoder(r)

	if err := decoder.Decode(to); err != nil {
		return fmt.Errorf("json decoder: %w", err)
	}

	return nil
}
