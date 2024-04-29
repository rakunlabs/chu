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

	var hooks []struct2.HookDecodeFunc
	if len(opt.Hooks) > 0 {
		hooks = make([]struct2.HookDecodeFunc, len(opt.Hooks))
		for i, h := range opt.Hooks {
			hooks[i] = struct2.HookDecodeFunc(h)
		}
	}

	decoder := struct2.Decoder{
		TagName:               loader.TagName,
		HooksDecode:           hooks,
		WeaklyTypedInput:      true,
		WeaklyIgnoreSeperator: opt.WeaklyIgnoreSeperator,
		WeaklyDashUnderscore:  opt.WeaklyDashUnderscore,
	}

	return &Loader{
		decoder: decoder,
	}
}

func (l Loader) SetValue(v interface{}) LoadSetter {
	l.value = v

	return LoadSetter(l)
}

// Map to load map to struct.
func (l LoadSetter) Load(_ context.Context, to any) error {
	if err := l.decoder.Decode(l.value, to); err != nil {
		return fmt.Errorf("mapx: %w", err)
	}

	return nil
}
