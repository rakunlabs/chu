package mapx

import (
	"context"
	"fmt"

	"github.com/rakunlabs/chu/loader"
	"github.com/worldline-go/struct2"
)

type Loader struct {
	value   interface{}
	decoder struct2.Decoder
}

type LoadSetter Loader

func New(opts ...Option) *Loader {
	opt := &option{
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  true,
		Tag:                   "cfg",
	}
	opt.apply(opts...)

	hooks := convertHookFuncs(opt.Hooks)

	return &Loader{
		decoder: struct2.Decoder{
			TagName:               opt.Tag,
			HooksDecode:           hooks,
			WeaklyTypedInput:      true,
			WeaklyIgnoreSeperator: opt.WeaklyIgnoreSeperator,
			WeaklyDashUnderscore:  opt.WeaklyDashUnderscore,
		},
	}
}

func (l Loader) SetValue(v interface{}) LoadSetter {
	l.value = v

	return LoadSetter(l)
}

func (l LoadSetter) Load(ctx context.Context, to any) error {
	return l.LoadChu(ctx, to)
}

// Map to load map to struct.
func (l LoadSetter) LoadChu(_ context.Context, to any, opts ...loader.Option) error {
	opt := loader.NewOption(opts...)

	hooks := convertHookFuncs(opt.Hooks)
	if len(hooks) > 0 {
		l.decoder.HooksDecode = hooks
	}

	if opt.Tag != "" {
		l.decoder.TagName = opt.Tag
	}

	if err := l.decoder.Decode(l.value, to); err != nil {
		return fmt.Errorf("mapx: %w", err)
	}

	return nil
}
