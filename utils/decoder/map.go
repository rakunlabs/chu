package decoder

import (
	"github.com/rakunlabs/chu/loader"
	"github.com/worldline-go/struct2"
)

type Map struct {
	decoder struct2.Decoder
}

func New(opts ...Option) *Map {
	opt := &option{
		WeaklyTypedInput:      true,
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  true,
		Tag:                   "cfg",
	}
	opt.apply(opts...)

	return &Map{
		decoder: struct2.Decoder{
			TagName:               opt.Tag,
			HooksDecode:           opt.Hooks,
			WeaklyTypedInput:      opt.WeaklyTypedInput,
			WeaklyIgnoreSeperator: opt.WeaklyIgnoreSeperator,
			WeaklyDashUnderscore:  opt.WeaklyDashUnderscore,
		},
	}
}

func (m *Map) Decode(input, output any) error {
	return m.decoder.Decode(input, output)
}

type Option func(*option)

type option struct {
	Hooks                 []loader.HookFunc
	WeaklyIgnoreSeperator bool
	WeaklyDashUnderscore  bool
	WeaklyTypedInput      bool
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
		o.Hooks = append(o.Hooks, hooks...)
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

// WithTag sets the tag for the configuration.
//   - default is "cfg"
func WithTag(tag string) Option {
	return func(o *option) {
		o.Tag = tag
	}
}

// WithWeaklyTypedInput sets the weakly typed input option.
//   - default is true
func WithWeaklyTypedInput(v bool) Option {
	return func(o *option) {
		o.WeaklyTypedInput = v
	}
}
