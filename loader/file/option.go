package file

import (
	"github.com/rakunlabs/chu/loader/mapx"
)

type Option func(*option)

type option struct {
	Mapx           *mapx.Loader
	FileSuffix     []string
	EtcFolderCheck bool
	Name           string
	Decoders       map[string]Decoder
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithMapxLoader(m *mapx.Loader) Option {
	return func(o *option) {
		o.Mapx = m
	}
}

func WithFileSuffix(suffix ...string) Option {
	return func(o *option) {
		o.FileSuffix = suffix
	}
}

func WithEtcFolderCheck(v bool) Option {
	return func(o *option) {
		o.EtcFolderCheck = v
	}
}

func WithDecoder(suffix string, d Decoder) Option {
	return func(o *option) {
		if o.Decoders == nil {
			o.Decoders = make(map[string]Decoder)
		}

		o.Decoders[suffix] = d
	}
}
