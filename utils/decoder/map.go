package decoder

import "github.com/rakunlabs/gofret"

type Map struct {
	codec *gofret.Codec
}

func New(opts ...Option) *Map {
	opt := &option{
		WeaklyTypedInput: true,
		Tag:              "cfg",
	}
	opt.apply(opts...)

	gopts := []gofret.Option{
		gofret.WithTag(opt.Tag),
	}
	gopts = append(gopts, opt.CodecOptions...)

	if opt.WeaklyTypedInput {
		gopts = append(gopts, gofret.WithWeakTypes())
	} else {
		gopts = append(gopts, gofret.WithStrictTypes())
	}

	if len(opt.Hooks) > 0 {
		gopts = append(gopts, gofret.WithHooks(opt.Hooks...))
	}

	return &Map{codec: gofret.New(gopts...)}
}

func (m *Map) Decode(input, output any) error {
	return m.codec.ToInto(input, output)
}

type Option func(*option)

type option struct {
	Hooks            []gofret.Hook
	CodecOptions     []gofret.Option
	WeaklyTypedInput bool
	Tag              string
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the gofret hooks for map to struct conversion.
func WithHooks(hooks ...gofret.Hook) Option {
	return func(o *option) {
		o.Hooks = append(o.Hooks, hooks...)
	}
}

// WithCodecOptions sets the gofret codec options.
func WithCodecOptions(opts ...gofret.Option) Option {
	return func(o *option) {
		o.CodecOptions = append(o.CodecOptions, opts...)
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
