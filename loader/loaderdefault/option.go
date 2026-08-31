package loaderdefault

import "github.com/rakunlabs/gofret"

type Option func(*option)

type option struct {
	Hooks        []gofret.Hook
	CodecOptions []gofret.Option
	TagName      string
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the gofret conversion hooks for the default loader.
func WithHooks(hooks ...gofret.Hook) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithCodecOptions sets the gofret codec options for the default loader.
func WithCodecOptions(opts ...gofret.Option) Option {
	return func(o *option) {
		o.CodecOptions = opts
	}
}

func WithTagName(tagName string) Option {
	return func(o *option) {
		o.TagName = tagName
	}
}
