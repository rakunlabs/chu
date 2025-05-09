package chu

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"

	// Enable default loaders.

	_ "github.com/rakunlabs/chu/loader/loaderdefault"
	_ "github.com/rakunlabs/chu/loader/loaderenv"
	_ "github.com/rakunlabs/chu/loader/loaderfile"
	_ "github.com/rakunlabs/chu/loader/loaderhttp"
)

var (
	DefaultLoaders = loader.Loaders
	DefaultHooks   = []loader.HookFunc{
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

	mapDecoder := decoder.New(
		decoder.WithTag(opt.Tag),
		decoder.WithHooks(opt.Hooks...),
		decoder.WithWeaklyIgnoreSeperator(opt.WeaklyIgnoreSeperator),
		decoder.WithWeaklyDashUnderscore(opt.WeaklyDashUnderscore),
	).Decode

	optLoader := loader.NewOption(
		loader.WithName(name),
		loader.WithHooks(opt.Hooks...),
		loader.WithTag(opt.Tag),
		loader.WithMapDecoder(mapDecoder),
		loader.WithLogger(opt.Logger),
	)

	loaderNames := loader.OrderLoaders(opt.Loaders)

	for _, lName := range loaderNames {
		l := opt.Loaders[lName]

		chuLoader := l.Loader()
		if err := chuLoader.LoadChu(ctx, to, optLoader); err != nil {
			if errors.Is(err, loader.ErrSkipLoader) {
				continue
			}

			return fmt.Errorf("config loader %s: %w", lName, err)
		}
	}

	return nil
}
