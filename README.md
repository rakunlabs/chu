# chu

[![License](https://img.shields.io/github/license/rakunlabs/chu?color=red&style=flat-square)](https://raw.githubusercontent.com/rakunlabs/chu/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/rakunlabs_chu?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=rakunlabs_chu)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/rakunlabs/chu/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/rakunlabs/chu/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/rakunlabs/chu?style=flat-square)](https://goreportcard.com/report/github.com/rakunlabs/chu)
[![Go PKG](https://raw.githubusercontent.com/rakunlabs/.github/main/assets/badges/gopkg.svg)](https://pkg.go.dev/github.com/rakunlabs/chu)

Configuration library to load from multiple sources.

```go
go get github.com/rakunlabs/chu
```

## Usage

Define a struct to hold the configuration.

```go
type Config struct {
    Name string   `cfg:"name"`
    Age  int      `cfg:"age"`
    Secret string `cfg:"secret" log:"-"` // skip this field in chu.Print
}
```

And load the configuration.

```go
cfg := Config{}

if err := chu.Load(ctx, "test", &cfg); err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

slog.Info("loaded configuration", "config", chu.Print(ctx, cfg))
```

The configuration will be loaded from the following sources in order:  
__-__ Default  
__-__ File  
__-__ Environment

`chu.Print` print the configuration in a human-readable format, skipping the fields with `log:"-"` or `log:"false"` tag. It uses `fmt.Stringer` interface to print the configuration.

### Loaders

Check [example](./example/) folder to see how to use loaders with different kind of configuration.

#### Default

Default loader is used to set default values from tag `default`.

```go
type Config struct {
    Name string `cfg:"name" default:"John"`
    Age  int    `cfg:"age"  default:"30"`
}
```

Default supports _numbers_, _string_, _bool_, _time.Duration_ and pointer of that types.

#### File

File loader is used to load configuration from file.

First checking `CONFIG_PATH` env value and try current location to find in order of `.toml`, `.yaml`, `.yml`, `.json` extension with using given name.

#### Environment

Environment loader is used to load configuration from environment variables.

`env` or `cfg` tag can usable for environment loader.

```go
export NAME=John
export AGE=30
```

```go
type Config struct {
    Name string `cfg:"name"`
    Age  int    `cfg:"age"`
}
```

When loading configuration, usable to change env loader's options.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loader.NameEnv, envloader.New(
        envloader.WithPrefix("MY_APP_"),
    )),
)
```

## Other Loaders

This loaders not enabled by default. You can use them by adding to the `DefaultLoaders` value.
