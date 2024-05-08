package chu

import (
	"github.com/rakunlabs/chu/loader"
)

type Option func(*option)

type option struct {
	Loaders []LoadHolder
	Hooks   []loader.HookFunc
	Tag     string
	// WeaklyIgnoreSeperator for map decoder option.
	//  - default is true
	//  - if true, ignore separator in map keys [-_ ]; "key1-key2" -> "key1key2"
	WeaklyIgnoreSeperator bool
	// WeaklyDashUnderscore for map decoder option.
	//  - default is false
	WeaklyDashUnderscore bool
	// Logger for logging.
	Logger loader.LogAdapter
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLoaders sets the loaders to use when loading the configuration.
//   - order matters
func WithLoaders(loaders ...LoadHolder) Option {
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

// WithLogger sets the logger for logging.
func WithLogger(logger loader.LogAdapter) Option {
	return func(o *option) {
		o.Logger = logger
	}
}
