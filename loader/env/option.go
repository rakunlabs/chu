package env

import "github.com/rakunlabs/chu/loader"

type Option func(*option)

type option struct {
	Hooks     []loader.HookFunc
	EnvHolder envHolder
	EnvFiles  []string
	TagName   string
	Prefix    string
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the hooks for the environment loader.
//   - if return type matches the field type, return value is assigned to the field
func WithHooks(hooks ...loader.HookFunc) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithEnv sets the environment variables to load and disable from the system environment.
func WithEnv(env map[string]string) Option {
	return func(o *option) {
		o.EnvHolder = envHolder(env)
	}
}

func WithEnvFile(path ...string) Option {
	return func(o *option) {
		o.EnvFiles = path
	}
}

func WithTagName(tagName string) Option {
	return func(o *option) {
		o.TagName = tagName
	}
}

// WithPrefix to just load the environment variables with the given prefix.
//   - if prefix is "APP_", then only the environment variables with the prefix "APP_" are loaded.
//   - prefix will strip from the field name.
func WithPrefix(prefix string) Option {
	return func(o *option) {
		o.Prefix = prefix
	}
}
