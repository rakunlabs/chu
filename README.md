# chu

Configuration library to load from multiple sources.

```go
go get github.com/rakunlabs/chu
```

## Usage

Define a struct to hold the configuration.

```go
type Config struct {
    Name string `cfg:"name"`
    Age  int    `cfg:"age"`
}
```

And load the configuration.

```go
cfg := Config{}

if err := chu.Load(ctx, "test", &cfg); err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

The configuration will be loaded from the following sources in order:  
__-__ Default  
__-__ File  
__-__ Environment

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

## Other Loaders

This loaders not enabled by default. You can use them by adding to the `Load` function.

