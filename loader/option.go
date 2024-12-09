package loader

import (
	"log/slog"

	"github.com/rakunlabs/logi"
)

type Option func(*option)

type option struct {
	Tag        string
	Name       string
	Hooks      []HookFunc
	MapDecoder func(input interface{}, output interface{}) error
	Logger     logi.Adapter
}

func NewOption(opts ...Option) *option {
	opt := &option{
		Name:   "",
		Tag:    "cfg",
		Logger: slog.Default(),
	}
	opt.apply(opts...)

	return opt
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMapDecoder sets the decoder for conversion between map and struct.
//   - output is the target struct
func WithMapDecoder(decoder func(input interface{}, output interface{}) error) Option {
	return func(o *option) {
		o.MapDecoder = decoder
	}
}

// WithName sets the name for loader.
//
// Loader will look this name for file, config name, etc.
func WithName(name string) Option {
	return func(o *option) {
		o.Name = name
	}
}

// WithHooks sets the hooks for conversion.
func WithHooks(hooks ...HookFunc) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithTag sets the tag name for struct field.
//   - loaders may use this tag to load the configuration.
func WithTag(tag string) Option {
	return func(o *option) {
		o.Tag = tag
	}
}

// WithLogger sets the logger for logging.
func WithLogger(logger logi.Adapter) Option {
	return func(o *option) {
		o.Logger = logger
	}
}
