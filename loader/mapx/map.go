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
	}
	opt.apply(opts...)

	hooks := convertHookFuncs(opt.Hooks)

	return &Loader{
		decoder: struct2.Decoder{
			TagName:               loader.TagName,
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

func (l LoadSetter) Load(ctx context.Context, ptr any) error {
	return l.LoadChu(ctx, ptr)
}

// Map to load map to struct.
func (l LoadSetter) LoadChu(_ context.Context, ptr any, opts ...loader.Option) error {
	opt := loader.NewOption(opts...)

	hooks := convertHookFuncs(opt.Hooks)
	if len(hooks) > 0 {
		l.decoder.HooksDecode = hooks
	}

	if err := l.decoder.Decode(l.value, ptr); err != nil {
		return fmt.Errorf("mapx: %w", err)
	}

	return nil
}
