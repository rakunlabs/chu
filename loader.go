package chu

import (
	"context"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/defaultx"
	"github.com/rakunlabs/chu/loader/env"
	"github.com/rakunlabs/chu/loader/file"
)

type Loader interface {
	LoadChu(ctx context.Context, ptr any, opts ...loader.Option) error
}

func Load(ctx context.Context, name string, ptr any, opts ...Option) error {
	opt := option{
		Loaders: []Loader{
			defaultx.New(),
			file.New(),
			env.New(),
		},
		Hooks: []loader.HookFunc{
			loader.HookTimeDuration,
		},
	}
	opt.apply(opts...)

	for _, l := range opt.Loaders {
		if err := l.LoadChu(
			ctx, ptr,
			loader.WithName(name),
			loader.WithHooks(opt.Hooks...),
		); err != nil {
			return err
		}
	}

	return nil
}
