package chu

import (
	"context"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/defaultx"
	"github.com/rakunlabs/chu/loader/env"
	"github.com/rakunlabs/chu/loader/file"
)

type Loader interface {
	Load(ctx context.Context, ptr any, opts ...loader.Option) error
}

var defaultLoaders = []Loader{
	defaultx.New(),
	file.New(),
	env.New(),
}

func Load(ctx context.Context, ptr any, opts ...Option) error {
	opt := option{
		Loaders: defaultLoaders,
	}
	opt.apply(opts...)

	for _, loader := range opt.Loaders {
		if err := loader.Load(ctx, ptr); err != nil {
			return err
		}
	}

	return nil
}
