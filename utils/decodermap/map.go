package decodermap

import (
	"github.com/rakunlabs/chu/loader"
	"github.com/worldline-go/struct2"
)

type Map struct {
	decoder struct2.Decoder
}

func New(opts ...Option) *Map {
	opt := &option{
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  false,
		Tag:                   "cfg",
	}
	opt.apply(opts...)

	hooks := convertHookFuncs(opt.Hooks)

	return &Map{
		decoder: struct2.Decoder{
			TagName:               opt.Tag,
			HooksDecode:           hooks,
			WeaklyTypedInput:      true,
			WeaklyIgnoreSeperator: opt.WeaklyIgnoreSeperator,
			WeaklyDashUnderscore:  opt.WeaklyDashUnderscore,
		},
	}
}

func (m *Map) Decode(input, output interface{}) error {
	return m.decoder.Decode(input, output)
}

func convertHookFuncs(hooks []loader.HookFunc) []struct2.HookDecodeFunc {
	if len(hooks) == 0 {
		return nil
	}

	hookFuncs := make([]struct2.HookDecodeFunc, len(hooks))
	for i, h := range hooks {
		hookFuncs[i] = struct2.HookDecodeFunc(h)
	}

	return hookFuncs
}

type Option func(*option)

type option struct {
	Hooks                 []loader.HookFunc
	WeaklyIgnoreSeperator bool
	WeaklyDashUnderscore  bool
	Tag                   string
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

// WithTag sets the tag for the configuration.
//   - default is "cfg"
func WithTag(tag string) Option {
	return func(o *option) {
		o.Tag = tag
	}
}
