package chu_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rakunlabs/chu"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/loaderenv"
)

func ExampleLoad() {
	cfg := struct {
		Name     string        `cfg:"name"`
		Age      int           `cfg:"age"`
		Duration time.Duration `cfg:"duration"`
	}{}

	_ = os.Setenv("MY_APP_NAME", "another")
	_ = os.Setenv("MY_APP_AGE", "70")
	_ = os.Setenv("MY_APP_DURATION", "1h30m")

	err := chu.Load(context.Background(), "my-app", &cfg, chu.WithLoaderOption(loader.NameEnv, loaderenv.New(loaderenv.WithPrefix("MY_APP_"))))
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Printf("%s\n", chu.MarshalJSON(cfg))
	// Output:
	// {"age":70,"duration":"1h30m0s","name":"another"}
}
