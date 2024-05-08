package consul

import (
	"io"
)

type Option func(*option)

type option struct {
	Decode func(r io.Reader, to interface{}) error
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithDecode sets the decoder for the consul loader.
//   - default is yaml decoder
func WithDecode(d func(r io.Reader, to interface{}) error) Option {
	return func(o *option) {
		o.Decode = d
	}
}
