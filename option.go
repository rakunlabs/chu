package chu

type Option func(*option)

type option struct {
	Loaders []Loader
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
