package fileloader

type Option func(*option)

type option struct {
	Folders    []string
	FileSuffix []string
	Decoders   map[string]Decoder
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithFolders sets the folders to use when loading the configuration.
//   - order matters
func WithFolders(folders ...string) Option {
	return func(o *option) {
		o.Folders = folders
	}
}

// WithFileSuffix sets the file suffixes to use when loading the configuration.
//   - order matters
func WithFileSuffix(suffixes ...string) Option {
	return func(o *option) {
		o.FileSuffix = suffixes
	}
}

// WithDecoders sets the decoders to use when loading the configuration.
//   - order matters
func WithDecoders(decoders map[string]Decoder) Option {
	return func(o *option) {
		o.Decoders = decoders
	}
}
