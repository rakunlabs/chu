package consul

import (
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/file"
)

type Option func(*option)

type option struct {
	Hooks                 []loader.HookFunc
	WeaklyIgnoreSeperator bool
	WeaklyDashUnderscore  bool
	Decoders              map[string]file.Decoder
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
//   - default is false
func WithWeaklyDashUnderscore(v bool) Option {
	return func(o *option) {
		o.WeaklyDashUnderscore = v
	}
}

func WithDecoder(suffix string, d file.Decoder) Option {
	return func(o *option) {
		if o.Decoders == nil {
			o.Decoders = make(map[string]file.Decoder)
		}

		o.Decoders[suffix] = d
	}
}
