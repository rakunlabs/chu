package defaultx

import "github.com/rakunlabs/chu/loader"

type Option func(*option)

type option struct {
	Hooks   []loader.HookFunc
	TagName string
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

func WithTagName(tagName string) Option {
	return func(o *option) {
		o.TagName = tagName
	}
}
