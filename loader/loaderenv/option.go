package loaderenv

import "github.com/rakunlabs/gofret"

type Option func(*option)

type option struct {
	Hooks              []gofret.Hook
	CodecOptions       []gofret.Option
	EnvHolder          envHolder
	CheckConfigEnvFile bool
	EnvFiles           []string
	TagEnv             string
	Tag                string
	Prefix             string
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the gofret conversion hooks for the environment loader.
func WithHooks(hooks ...gofret.Hook) Option {
	return func(o *option) {
		o.Hooks = hooks
	}
}

// WithCodecOptions sets the gofret codec options for the environment loader.
func WithCodecOptions(opts ...gofret.Option) Option {
	return func(o *option) {
		o.CodecOptions = opts
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

// WithTagEnv is `env` by default.
func WithTagEnv(tagName string) Option {
	return func(o *option) {
		o.TagEnv = tagName
	}
}

// WithTag is `cfg` by default.
//   - tag is used if `env“ tag is not found.
func WithTag(tag string) Option {
	return func(o *option) {
		o.Tag = tag
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

// WithCheckConfigEnvFile enables checking CONFIG_ENV_FILE variable to load specific environment file.
//   - if true and CONFIG_ENV_FILE is set, the file is appended to the list of env files to load.
//   - default is true.
func WithCheckConfigEnvFile(check bool) Option {
	return func(o *option) {
		o.CheckConfigEnvFile = check
	}
}
