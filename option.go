package chu

import (
	"github.com/rakunlabs/logi/logadapter"

	"github.com/rakunlabs/chu/loader"
)

type Option func(*option)

type option struct {
	Loaders map[loader.LoaderName]loader.Loader
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
	Logger logadapter.Adapter
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithLoader adds a loader to the configuration.
func WithLoader(ld loader.Loader) Option {
	return func(o *option) {
		o.Loaders[ld.LoadName()] = ld
	}
}

// WithLoaderOption sets the loader option for the configuration.
//   - ld is the loader to set.
//   - if the loader name is not exist, it will be ignored
func WithLoaderOption(ld loader.Loader) Option {
	return func(o *option) {
		name := ld.LoadName()

		if _, ok := o.Loaders[name]; ok {
			o.Loaders[name] = ld
		}
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
func WithLogger(logger logadapter.Adapter) Option {
	return func(o *option) {
		o.Logger = logger
	}
}

// WithLoaderNames sets the loaders inside default loaders use when loading the configuration.
func WithLoaderNames(names ...loader.LoaderName) Option {
	return func(o *option) {
		loaders := make(map[loader.LoaderName]loader.Loader, len(names))
		for _, name := range names {
			if _, ok := o.Loaders[name]; !ok {
				continue
			}
			loaders[name] = o.Loaders[name]
		}

		o.Loaders = loaders
	}
}

// WithDisableLoader disables a loader inside default loaders use when loading the configuration.
func WithDisableLoader(names ...loader.LoaderName) Option {
	return func(o *option) {
		for _, name := range names {
			delete(o.Loaders, name)
		}
	}
}
