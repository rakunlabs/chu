package chu

import "github.com/rakunlabs/chu/loader"

type Option func(*option)

type option struct {
	Loaders []Loader
	Hooks   []loader.HookFunc
	Tag     string
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLoaders sets the loaders to use when loading the configuration.
//   - order matters
func WithLoaders(loaders ...Loader) Option {
	return func(o *option) {
		o.Loaders = loaders
	}
}

// WithHookSet sets the hooks for conversion.
func WithHookSet(hooks ...loader.HookFunc) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithHook adds hooks for conversion.
func WithHook(hooks ...loader.HookFunc) Option {
	return func(o *option) {
		o.Hooks = append(o.Hooks, hooks...)
	}
}

// WithTag sets the tag for the configuration.
//   - default is "cfg"
func WithTag(tag string) Option {
	return func(o *option) {
		o.Tag = tag
	}
}
