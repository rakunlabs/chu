package chu_test

import (
	"context"
	"fmt"
	"os"

	"github.com/rakunlabs/chu"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/envloader"
)

func ExampleLoad() {
	cfg := struct {
		Name string `cfg:"name"`
		Age  int    `cfg:"age"`
	}{}

	_ = os.Setenv("MY_APP_NAME", "another")
	_ = os.Setenv("MY_APP_AGE", "70")

	err := chu.Load(context.Background(), "my-app", &cfg, chu.WithLoaderOption(loader.NameEnv, envloader.New(envloader.WithPrefix("MY_APP_"))))
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Printf("%s\n", chu.MarshalJSON(cfg))
	// Output:
	// {"age":70,"name":"another"}
}
