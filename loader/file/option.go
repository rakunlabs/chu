package file

import (
	"github.com/rakunlabs/chu/loader"
)

type Option func(*option)

type option struct {
	Hooks                 []loader.HookFunc
	WeaklyIgnoreSeperator bool
	WeaklyDashUnderscore  bool
	FileSuffix            []string
	EtcFolderCheck        bool
	Name                  string
	Decoders              map[string]Decoder
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the hooks for map to struct conversion.
func WithHooks(hooks ...loader.HookFunc) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithWeaklyIgnoreSeperator sets the weakly ignore separator option.
//   - default is true
func WithWeaklyIgnoreSeperator(v bool) Option {
	return func(o *option) {
		o.WeaklyIgnoreSeperator = v
	}
}

// WithWeaklyDashUnderscore sets the weakly dash underscore option.
//   - default is true
func WithWeaklyDashUnderscore(v bool) Option {
	return func(o *option) {
		o.WeaklyDashUnderscore = v
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
