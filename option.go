package chu

import (
	"github.com/rakunlabs/gofret"
	"github.com/rakunlabs/logi/logadapter"

	"github.com/rakunlabs/chu/loader"
)

type Option func(*option)

type option struct {
	Loaders      map[loader.LoaderName]loader.Loader
	Hooks        []gofret.Hook
	CodecOptions []gofret.Option
	Tag          string
	// Logger for logging.
	Logger logadapter.Adapter
	// Version of the application, used by some loaders (e.g. http loader as `version` query).
	Version string
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

// WithHookSet sets the gofret hooks for conversion.
func WithHookSet(hooks ...gofret.Hook) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithHook adds gofret hooks for conversion.
func WithHook(hooks ...gofret.Hook) Option {
	return func(o *option) {
		o.Hooks = append(o.Hooks, hooks...)
	}
}

// WithCodecOptions adds gofret codec options for conversion.
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

// WithLogger sets the logger for logging.
func WithLogger(logger logadapter.Adapter) Option {
	return func(o *option) {
		o.Logger = logger
	}
}

// WithVersion sets the version of the application.
//   - Loaders may use it to fetch a specific configuration version.
//   - HTTP loader appends it as `version=<version>` query when no `CONFIG_HTTP_QUERY` env is set.
func WithVersion(version string) Option {
	return func(o *option) {
		o.Version = version
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
