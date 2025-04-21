package basic

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/rakunlabs/chu"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/envloader"
)

type Config struct {
	Host    string `cfg:"host"     default:"localhost"`
	Port    int    `cfg:"port"     default:"8080"`
	PortPtr *int   `cfg:"port_ptr" default:"8080"`
	Test    int    `cfg:"test"     default:"1"`

	Duration    time.Duration  `cfg:"duration"     default:"1s"`
	DurationPtr *time.Duration `cfg:"duration_ptr" default:"2s"`

	// Database configuration
	DB struct {
		Pass string `cfg:"pass"` // DB_PASS environment variable
	}

	Fn      func()     `log:"false"` // cannot be loaded, result is <nil>
	Channel <-chan int // cannot be loaded, result is <nil>

	// Special configuration
	Special SpecialConfig `cfg:"special"`
}

type SpecialConfig struct {
	Host string `cfg:"host" default:"localhost"`
	Port string `cfg:"port" default:"8080"`
}

func (c *SpecialConfig) String() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Load(ctx context.Context) error {
	cfg := Config{}

	_ = os.Setenv("MY_APP_DB_PASS", "password")
	_ = os.Setenv("CONFIG_FILE", "basic/testdata/app.toml")

	if err := chu.Load(ctx, "my-app", &cfg,
		chu.WithLoaderOption(loader.NameEnv, envloader.New(
			envloader.WithPrefix("MY_APP_"),
		)),
	); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	slog.Info("loaded configuration", "config", chu.Print(ctx, cfg))

	return nil
}
