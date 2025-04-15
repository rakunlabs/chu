package chu

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/defaultloader"
	"github.com/rakunlabs/chu/loader/envloader"
	"github.com/rakunlabs/chu/loader/fileloader"
	"github.com/rakunlabs/chu/utils/decodermap"
)

type LoadHolder struct {
	Name   string
	Loader func() loader.Loader
	Order  *Order
}

type Order struct {
	Before string
	After  string
}

var (
	DefaultLoaders = []LoadHolder{
		{Name: defaultloader.LoaderName, Loader: defaultloader.New()},
		{Name: fileloader.LoaderName, Loader: fileloader.New()},
		{Name: envloader.LoaderName, Loader: envloader.New()},
	}
	DefaultHooks = []loader.HookFunc{
		loader.HookTimeDuration,
	}
	DefaultOptions = []Option{}
)

// Load loads the configuration from loaders.
//   - default loaders are [defaultx, file, env].
//   - default hooks are [loader.HookTimeDuration].
func Load(ctx context.Context, name string, to any, opts ...Option) error {
	opts = append(DefaultOptions, opts...)

	opt := option{
		Loaders:               DefaultLoaders,
		Hooks:                 DefaultHooks,
		Tag:                   "cfg",
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  false,
		Logger:                slog.Default(),
	}
	opt.apply(opts...)

	mapDecoder := decodermap.New(
		decodermap.WithTag(opt.Tag),
		decodermap.WithHooks(opt.Hooks...),
		decodermap.WithWeaklyIgnoreSeperator(opt.WeaklyIgnoreSeperator),
		decodermap.WithWeaklyDashUnderscore(opt.WeaklyDashUnderscore),
	).Decode

	optLoader := loader.NewOption(
		loader.WithName(name),
		loader.WithHooks(opt.Hooks...),
		loader.WithTag(opt.Tag),
		loader.WithMapDecoder(mapDecoder),
		loader.WithLogger(opt.Logger),
	)

	for _, l := range opt.Loaders {
		chuLoader := l.Loader()
		if err := chuLoader.LoadChu(ctx, to, optLoader); err != nil {
			if errors.Is(err, loader.ErrSkipLoader) {
				opt.Logger.Debug(err.Error(), "loader", l.Name)

				continue
			}

			return fmt.Errorf("config loader %s: %w", l.Name, err)
		}
	}

	return nil
}
