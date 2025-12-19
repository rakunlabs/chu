package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"strings"

	"github.com/rakunlabs/chu/example/basic"
	"github.com/rakunlabs/chu/example/mix"
)

var Examples = map[string]Exampler{
	"1": {
		Name: "basic struct configuration",
		Fn:   basic.Load,
	},
	"2": {
		Name: "vault and consul configuration",
		Fn:   mix.Load,
	},
}

type Exampler struct {
	Name string
	Fn   func(context.Context) error
}

func main() {
	n := getNumber()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("EXAMPLE_NO: " + n)

	v, ok := Examples[n]
	if !ok {
		slog.Error("invalid example number")

		return
	}

	slog.Info("running example: " + v.Name)
	if err := v.Fn(context.Background()); err != nil {
		slog.Error("example", "error", err)

		os.Exit(1)
	}
}

func getNumber() string {
	var (
		n    string
		help bool
	)

	flag.StringVar(&n, "number", "1", "example number")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()

	if help {
		flag.PrintDefaults()

		os.Exit(0)
	}

	if n == "" {
		if v := strings.TrimSpace(os.Getenv("EXAMPLE_NO")); v != "" {
			n = v
		}
	}

	if n == "" {
		n = "1"
	}

	return n
}
