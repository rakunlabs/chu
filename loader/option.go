package loader

import (
	"log/slog"
	"os"
	"strings"

	"github.com/rakunlabs/logi/logadapter"
)

type OptionFunc func(*Option)

type Option struct {
	Tag        string
	Name       string
	Hooks      []HookFunc
	MapDecoder func(input any, output any) error
	Logger     logadapter.Adapter
}

func NewOption(opts ...OptionFunc) *Option {
	opt := &Option{
		Name:   "",
		Tag:    "cfg",
		Logger: slog.Default(),
	}
	opt.apply(opts...)

	return opt
}

func (o *Option) apply(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(o)
	}

	if v := os.Getenv("CONFIG_NAME_PREFIX"); v != "" {
		o.Name = strings.Trim(strings.Trim(v, "/")+"/"+strings.Trim(o.Name, "/"), "/")
	}
}

// WithMapDecoder sets the decoder for conversion between map and struct.
//   - output is the target struct
func WithMapDecoder(decoder func(input any, output any) error) OptionFunc {
	return func(o *Option) {
		o.MapDecoder = decoder
	}
}

// WithName sets the name for loader.
//
// Loader will look this name for file, config name, etc.
func WithName(name string) OptionFunc {
	return func(o *Option) {
		o.Name = name
	}
}

// WithHooks sets the hooks for conversion.
func WithHooks(hooks ...HookFunc) OptionFunc {
	return func(o *Option) {
		o.Hooks = hooks
	}
}

// WithTag sets the tag name for struct field.
//   - loaders may use this tag to load the configuration.
func WithTag(tag string) OptionFunc {
	return func(o *Option) {
		o.Tag = tag
	}
}

// WithLogger sets the logger for logging.
func WithLogger(logger logadapter.Adapter) OptionFunc {
	return func(o *Option) {
		o.Logger = logger
	}
}
