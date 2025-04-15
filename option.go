package chu

import (
	"github.com/rakunlabs/logi/logadapter"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/defaultloader"
	"github.com/rakunlabs/chu/loader/envloader"
	"github.com/rakunlabs/chu/loader/fileloader"
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
	Logger logadapter.Adapter
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

func WithEnvLoaderOptions(opts ...envloader.Option) Option {
	return func(o *option) {
		for i, l := range o.Loaders {
			if l.Name == envloader.LoaderName {
				o.Loaders[i].Loader = envloader.New(opts...)

				break
			}
		}
	}
}

func WithDefaultLoaderOptions(opts ...defaultloader.Option) Option {
	return func(o *option) {
		for i, l := range o.Loaders {
			if l.Name == defaultloader.LoaderName {
				o.Loaders[i].Loader = defaultloader.New(opts...)

				break
			}
		}
	}
}

// WithFileLoaderOptions sets the file loader options.
func WithFileLoaderOptions(opts ...fileloader.Option) Option {
	return func(o *option) {
		for i, l := range o.Loaders {
			if l.Name == fileloader.LoaderName {
				o.Loaders[i].Loader = fileloader.New(opts...)

				break
			}
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
func WithLoaderNames(names ...string) Option {
	return func(o *option) {
		loaders := make([]LoadHolder, 0, len(names))
		for _, name := range names {
			for _, l := range o.Loaders {
				if l.Name == name {
					loaders = append(loaders, l)

					break
				}
			}
		}

		o.Loaders = loaders
	}
}
