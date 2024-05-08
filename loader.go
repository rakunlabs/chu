package chu

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/consul"
	"github.com/rakunlabs/chu/loader/defaultx"
	"github.com/rakunlabs/chu/loader/env"
	"github.com/rakunlabs/chu/loader/file"
	"github.com/rakunlabs/chu/loader/vault"
	"github.com/rakunlabs/chu/utils/decoder"
)

type Loader interface {
	LoadChu(ctx context.Context, to any, opts ...loader.Option) error
}

type LoadHolder struct {
	Name   string
	Loader Loader
}

var (
	defaultLoaders = []LoadHolder{
		{Name: "default", Loader: defaultx.New()},
		{Name: "consul", Loader: consul.New()},
		{Name: "vault", Loader: vault.New()},
		{Name: "file", Loader: file.New()},
		{Name: "env", Loader: env.New()},
	}
	defaultHooks = []loader.HookFunc{
		loader.HookTimeDuration,
	}
)

// Load loads the configuration from loaders.
//   - default loaders are [defaultx, file, env].
//   - default hooks are [loader.HookTimeDuration].
func Load(ctx context.Context, name string, to any, opts ...Option) error {
	opt := option{
		Loaders:               defaultLoaders,
		Hooks:                 defaultHooks,
		Tag:                   "cfg",
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  false,
		Logger:                slog.Default(),
	}
	opt.apply(opts...)

	mapDecoder := decoder.NewMap(
		decoder.WithTag(opt.Tag),
		decoder.WithHooks(opt.Hooks...),
		decoder.WithWeaklyIgnoreSeperator(opt.WeaklyIgnoreSeperator),
		decoder.WithWeaklyDashUnderscore(opt.WeaklyDashUnderscore),
	).Decode

	for _, l := range opt.Loaders {
		if err := l.Loader.LoadChu(
			ctx, to,
			loader.WithName(name),
			loader.WithHooks(opt.Hooks...),
			loader.WithTag(opt.Tag),
			loader.WithMapDecoder(mapDecoder),
			loader.WithLogger(opt.Logger),
		); err != nil {
			if errors.Is(err, loader.ErrSkipLoader) {
				opt.Logger.Debug(err.Error(), "loader", l.Name)

				continue
			}

			return fmt.Errorf("config loader %s: %w", l.Name, err)
		}
	}

	return nil
}
