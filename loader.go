package chu

import (
	"context"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/defaultx"
	"github.com/rakunlabs/chu/loader/env"
	"github.com/rakunlabs/chu/loader/file"
)

type Loader interface {
	LoadChu(ctx context.Context, to any, opts ...loader.Option) error
}

// Load loads the configuration from loaders.
//   - default loaders are [defaultx, file, env].
//   - default hooks are [loader.HookTimeDuration].
func Load(ctx context.Context, name string, to any, opts ...Option) error {
	opt := option{
		Loaders: []Loader{
			defaultx.New(),
			file.New(),
			env.New(),
		},
		Hooks: []loader.HookFunc{
			loader.HookTimeDuration,
		},
		Tag: "cfg",
	}
	opt.apply(opts...)

	for _, l := range opt.Loaders {
		if err := l.LoadChu(
			ctx, to,
			loader.WithName(name),
			loader.WithHooks(opt.Hooks...),
			loader.WithTag(opt.Tag),
		); err != nil {
			return err
		}
	}

	return nil
}
