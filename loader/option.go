package loader

type Option func(*option)

type option struct {
	Tag   string
	Name  string
	Hooks []HookFunc
}

func NewOption(opts ...Option) *option {
	opt := &option{
		Name: "",
		Tag:  "cfg",
	}
	opt.apply(opts...)
	return opt
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
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
