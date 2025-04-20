package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/rakunlabs/chu/example/basic"
)

var Examples = map[int]Exampler{
	1: {
		Name: "basic struct configuration",
		Fn:   basic.Load,
	},
}

type Exampler struct {
	Name string
	Fn   func(context.Context)
}

func main() {
	n := getNumber()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info(fmt.Sprintf("EXAMPLE_NO: %d", n))

	v, ok := Examples[n]
	if !ok {
		slog.Error("invalid example number")

		return
	}

	slog.Info(fmt.Sprintf("running example: %s", v.Name))
	v.Fn(context.Background())
}

func getNumber() int {
	var (
		n    int
		help bool
	)

	flag.IntVar(&n, "number", 1, "example number")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()

	if help {
		flag.PrintDefaults()

		os.Exit(0)
	}

	if n == 0 {
		if v := strings.TrimSpace(os.Getenv("EXAMPLE_NO")); v != "" {
			n, _ = strconv.Atoi(v)
		}
	}

	if n == 0 {
		n = 1
	}

	return n
}
