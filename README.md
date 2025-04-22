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
    Secret string `cfg:"secret" log:"-"` // skip this field in chu.MarshalMap
}
```

And load the configuration.

```go
cfg := Config{}

if err := chu.Load(ctx, "test", &cfg); err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

slog.Info("loaded configuration", "config", chu.MarshalMap(ctx, cfg))
```

The configuration will be loaded from the following sources in order:  
__-__ Default  
__-__ File  
__-__ Http  
__-__ Environment

`chu.MarshalMap` or `chu.MarshalJSON` print the configuration, skipping the fields `log:"false"` tag and value unless `1, t, T, TRUE, true, True` makes false.  
String func use `fmt.Stringer` interface checks to print the configuration.

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

#### Http

HTTP loader is used to load configuration from HTTP server.

| Env Value            | Description                                      | Default |
| -------------------- | ------------------------------------------------ | ------- |
| `CONFIG_HTTP_ADDR`   | HTTP server address, not exist than skips loader | -       |
| `CONFIG_HTTP_PREFIX` | Prefix for the configuration                     | -       |

It send `GET` request to the server with `CONFIG_HTTP_ADDR` env value with appending the name as path.  
`204` or `404` response code will skip the loader, only accept `200` response code.

#### File

File loader is used to load configuration from file.

First checking `CONFIG_FILE` env value and try current location to find in order of `.toml`, `.yaml`, `.yml`, `.json` extension with using given name.

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

### Other Loaders

This loaders not enabled by default. Import the package to enable it.  
Use `chu.WithLoaderOption` to set the loader options.  
Or use `chu.WithLoader` to set the loaders manually.

<details><summary>#### Vault</summary>

Vault loader is used to load configuration from HashiCorp Vault.  
This is not enabled by default.

Enable Vault loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/vaultloader"
)
```

| Env Value                       | Description                                          | Default              |
| ------------------------------- | ---------------------------------------------------- | -------------------- |
| `VAULT_SECRET_BASE_PATH`        | Prefix for the configuration, must given base        | -                    |
| `VAULT_ADDR` `VAULT_AGENT_ADDR` | Vault server address, not exist than skips loader    | -                    |
| `VAULT_ROLE_ID`                 | Role ID for AppRole authentication, for role login   | -                    |
| `VAULT_SECRET_ID`               | Secret ID for AppRole authentication, for role login | -                    |
| `VAULT_APPROLE_BASE_PATH`       | Base path for AppRole authentication, for role login | `auth/approle/login` |

</details>

<details><summary>#### Consul</summary>

Consul loader is used to load configuration from HashiCorp Consul.  
This is not enabled by default.

Enable Consul loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/consulloader"
)
```

| Env Value                   | Description                                        | Default |
| --------------------------- | -------------------------------------------------- | ------- |
| `CONSUL_CONFIG_PATH_PREFIX` | Prefix for the configuration                       | -       |
| `CONSUL_HTTP_ADDR`          | Consul server address, not exist than skips loader | -       |

</details>
