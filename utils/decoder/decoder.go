package decoder

import (
	"io"
)

type Decoder interface {
	Decode(r io.Reader, to interface{}) error
}
