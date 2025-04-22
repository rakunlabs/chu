package mix

import (
	"context"
	"log/slog"

	"github.com/rakunlabs/chu"

	_ "github.com/rakunlabs/chu/loader/consulloader"
	_ "github.com/rakunlabs/chu/loader/vaultloader"
)

type Config struct {
	Test int `cfg:"test" default:"1"`

	// Database configuration
	DB struct {
		Pass string `cfg:"pass" log:"-"` // DB_PASS environment variable
	} `cfg:"db"`
}

func Load(ctx context.Context) error {
	cfg := Config{}

	if err := chu.Load(ctx, "app/mix", &cfg); err != nil {
		return err
	}

	slog.Info("loaded configuration", "config", chu.MarshalJSON(cfg))

	return nil
}
