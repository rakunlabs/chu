package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/rakunlabs/chu"
)

var MapExample = map[string]Exampler{
	"1": {
		Name: "basic struct configuration",
		Fn:   basicStruct,
	},
}

type Exampler struct {
	Name string
	Fn   func()
}

type Config struct {
	Host string `cfg:"host" default:"localhost"`
	Port int    `cfg:"port" default:"8080"`

	// Database configuration
	DB struct {
		Pass string `cfg:"pass"`
	}
}

func basicStruct() {
	cfg := Config{}

	_ = os.Setenv("DB_PASS", "password")

	if err := chu.Load(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %+v", cfg)
}

func main() {
	var n string
	flag.StringVar(&n, "number", "", "example number")
	flag.Parse()

	if n == "" {
		n = os.Getenv("EXAMPLE_NO")
	}

	if n == "" {
		n = "1"
	}

	log.Printf("EXAMPLE_NO: %s", n)

	v, ok := MapExample[n]
	if !ok {
		log.Println("invalid example number")

		return
	}

	log.Printf("running example: %s", v.Name)
	v.Fn()
}
