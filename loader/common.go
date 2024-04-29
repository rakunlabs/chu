package loader

import "reflect"

// HookFunc get input, output and data and return modified data.
type HookFunc func(reflect.Type, reflect.Type, interface{}) (interface{}, error)

type Option func(*option)

type option struct {
	Name string
}

func NewOption(opts ...Option) *option {
	opt := &option{
		Name: "",
	}
	opt.apply(opts...)
	return opt
}

func (o *option) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithHooks sets the hooks for map to struct conversion.
func WithName(name string) Option {
	return func(o *option) {
		o.Name = name
	}
}
