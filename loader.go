package chu

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
	"github.com/rakunlabs/gofret"
	"github.com/xhit/go-str2duration/v2"

	// Enable default loaders.

	_ "github.com/rakunlabs/chu/loader/loaderdefault"
	_ "github.com/rakunlabs/chu/loader/loaderenv"
	_ "github.com/rakunlabs/chu/loader/loaderfile"
	_ "github.com/rakunlabs/chu/loader/loaderhttp"
)

var (
	DefaultLoaders      = loader.Loaders
	DefaultHooks        = []gofret.Hook{}
	DefaultCodecOptions = []gofret.Option{
		gofret.Between(str2duration.ParseDuration),
	}
	DefaultOptions = []Option{}
)

// Load loads the configuration from loaders.
//   - default loaders are [defaultx, file, env].
func Load(ctx context.Context, name string, to any, opts ...Option) error {
	opts = append(DefaultOptions, opts...)

	opt := option{
		Loaders:      DefaultLoaders,
		Hooks:        DefaultHooks,
		CodecOptions: append([]gofret.Option{}, DefaultCodecOptions...),
		Tag:          "cfg",
		Logger:       slog.Default(),
	}
	opt.apply(opts...)

	mapDecoder := decoder.New(
		decoder.WithTag(opt.Tag),
		decoder.WithHooks(opt.Hooks...),
		decoder.WithCodecOptions(opt.CodecOptions...),
	).Decode

	optLoader := loader.NewOption(
		loader.WithName(name),
		loader.WithHooks(opt.Hooks...),
		loader.WithCodecOptions(opt.CodecOptions...),
		loader.WithTag(opt.Tag),
		loader.WithMapDecoder(mapDecoder),
		loader.WithLogger(opt.Logger),
		loader.WithVersion(opt.Version),
	)

	loaderNames := loader.OrderLoaders(opt.Loaders)

	for _, lName := range loaderNames {
		if v, ok := os.LookupEnv("CONFIG_SET_" + strings.ToUpper(string(lName))); ok {
			if enabled, err := strconv.ParseBool(v); err == nil && !enabled {
				optLoader.Logger.Info("config loader disabled via env", "loader", lName)

				continue
			}
		}

		l := opt.Loaders[lName]

		if err := l.Load(ctx, to, optLoader); err != nil {
			if errors.Is(err, loader.ErrSkipLoader) {
				continue
			}

			return fmt.Errorf("config loader %s: %w", lName, err)
		}
	}

	return nil
}
