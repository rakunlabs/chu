package basic

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rakunlabs/chu"
)

type Config struct {
	Host    string `cfg:"host"     default:"localhost"`
	Port    int    `cfg:"port"     default:"8080"`
	PortPtr *int   `cfg:"port_ptr" default:"8080"`

	Duration    time.Duration  `cfg:"duration"     default:"1s"`
	DurationPtr *time.Duration `cfg:"duration_ptr" default:"2s"`

	// Database configuration
	DB struct {
		Pass string `cfg:"pass"` // DB_PASS environment variable
	}

	Fn      func()     // cannot be loaded, result is <nil>
	Channel <-chan int // cannot be loaded, result is <nil>
}

func Load(ctx context.Context) {
	cfg := Config{}

	_ = os.Setenv("DB_PASS", "password")

	if err := chu.Load(ctx, "test", &cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v", cfg)
}
