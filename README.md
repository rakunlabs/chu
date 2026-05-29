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

`CONFIG_NAME_PREFIX` env value can be used to set the prefix for the configuration.

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

| Env Value            | Description                                                              | Default |
| -------------------- | ------------------------------------------------------------------------ | ------- |
| `CONFIG_HTTP_ADDR`   | HTTP server address, not exist than skips loader                         | -       |
| `CONFIG_HTTP_SUFFIX` | Suffix for the configuration path                                        | -       |
| `CONFIG_HTTP_QUERY`  | Query parameters appended to the request URL (e.g. `env=prod&region=us`) | -       |

It send `GET` request to the server with `CONFIG_HTTP_ADDR` env value with appending the name as path.  
`204` or `404` response code will skip the loader, only accept `200` response code.

When `CONFIG_HTTP_QUERY` is not set, the loader falls back to `chu.WithVersion` and appends `?version=<version>` to the request URL.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithVersion("1.2.3"),
)
```

#### File

File loader is used to load configuration from file.

First checking `CONFIG_FILE` env value and try current location to find in order of `.toml`, `.yaml`, `.yml`, `.json` extension with using given name.

#### Environment

Environment loader is used to load configuration from environment variables.

`env` or `cfg` tag can usable for environment loader.  
If you have `.env`, `.env.local` file in the current directory, it will be loaded automatically.  
Add extra env file using `CONFIG_ENV_FILE` env value.

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
    // now you need to set prefix "MY_APP_" for env variables
    chu.WithLoaderOption(loaderenv.New(
        loaderenv.WithPrefix("MY_APP_"),
    )),
)
```

Use the `noprefix` tag option on fields that should also check unprefixed env vars as a fallback.  
The prefixed value takes priority when both exist.

```go
type Config struct {
    AppName  string `cfg:"app_name"`              // only MY_APP_APP_NAME
    LogLevel string `cfg:"log_level,noprefix"`    // MY_APP_LOG_LEVEL, falls back to LOG_LEVEL
}
```

#### Loader toggles

Every registered loader can be skipped with `CONFIG_SET_<LOADER_NAME>=false`.
This is useful when an optional external backend is down and should not block
startup. Unset, `true`, or invalid values keep the loader enabled.

```sh
CONFIG_SET_VAULT=false
CONFIG_SET_CONSUL=false
CONFIG_SET_AWSSECRETS=false
```

Loader names are uppercased from their registered names, for example:
`CONFIG_SET_DEFAULT`, `CONFIG_SET_CONSUL`, `CONFIG_SET_AWSSSM`,
`CONFIG_SET_AWSSECRETS`, `CONFIG_SET_AZUREKEYVAULT`, `CONFIG_SET_VAULT`,
`CONFIG_SET_GCPPARAMETER`, `CONFIG_SET_GCPSECRET`, `CONFIG_SET_HTTP`,
`CONFIG_SET_FILE`, `CONFIG_SET_ENV`.

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
    _ "github.com/rakunlabs/chu/loader/external/loadervault"
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
    _ "github.com/rakunlabs/chu/loader/external/loaderconsul"
)
```

| Env Value                   | Description                                        | Default |
| --------------------------- | -------------------------------------------------- | ------- |
| `CONSUL_CONFIG_PATH_PREFIX` | Prefix for the configuration                       | -       |
| `CONSUL_HTTP_ADDR`          | Consul server address, not exist than skips loader | -       |

</details>

<details><summary>#### GCP Parameter Manager</summary>

GCP Parameter Manager loader is used to load configuration from Google Cloud Parameter Manager.  
This is not enabled by default.

Enable GCP Parameter Manager loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/loader/external/loadergcpparameter"
)
```

| Env Value                | Description                                                                                         | Default  |
| ------------------------ | --------------------------------------------------------------------------------------------------- | -------- |
| `GCP_PROJECT_ID`         | GCP project id, not exist than skips loader (env > option)                                          | -        |
| `GCP_PARAMETER_LOCATION` | Parameter Manager location                                                                          | `global` |
| `GCP_PARAMETER_VERSION`  | Parameter version to render                                                                         | `latest` |
| `GCP_PARAMETER_FALLBACK` | When `latest` is disabled, fall back to the most recent enabled version (env > option), e.g. `true` | `false`  |

It calls `RenderParameterVersion` so referenced secrets are resolved server-side.  
The configuration name is used as the parameter id; `/` characters in the name are replaced with `-` because GCP parameter ids do not allow `/`.  
Authentication uses Application Default Credentials (`GOOGLE_APPLICATION_CREDENTIALS`, workload identity, etc.).  
Payload is decoded as YAML by default (JSON is a valid YAML subset).

Project id can also be injected programmatically. The `GCP_PROJECT_ID` env var, when set, takes priority over the option.

By default, when `GCP_PARAMETER_VERSION` is unset (`latest`) and that latest version is disabled, the render fails. Enable version fallback to instead render the most recently created enabled version (selected by create time, since parameter version ids are arbitrary strings). Fallback only applies to `latest`; an explicit version is never overridden. The `GCP_PARAMETER_FALLBACK` env var, when set, takes priority over the option.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loadergcpparameter.New(
        loadergcpparameter.WithProjectID("my-gcp-project"),
        loadergcpparameter.WithVersionFallback(true),
    )),
)
```

</details>

<details><summary>#### GCP Secret Manager</summary>

GCP Secret Manager loader is used to load configuration from Google Cloud Secret Manager.  
This is not enabled by default.

Enable GCP Secret Manager loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/loader/external/loadergcpsecret"
)
```

| Env Value               | Description                                                                                                    | Default  |
| ----------------------- | -------------------------------------------------------------------------------------------------------------- | -------- |
| `GCP_PROJECT_ID`        | GCP project id, not exist than skips loader (env > option)                                                     | -        |
| `GCP_SECRET_LOCATION`   | Secret Manager location (use `global` or a region)                                                             | `global` |
| `GCP_SECRET_VERSION`    | Secret version to access                                                                                       | `latest` |
| `GCP_SECRET_ADDITIONAL` | Comma-separated list of shared secret ids loaded BEFORE the app secret (env > option), e.g. `generic,keycloak` | -        |

It calls `AccessSecretVersion` to fetch the secret payload.  
The configuration name is used as the secret id; `/` characters in the name are replaced with `-` because GCP secret ids do not allow `/`.  
Authentication uses Application Default Credentials (`GOOGLE_APPLICATION_CREDENTIALS`, workload identity, etc.).  
Payload is decoded as YAML by default (JSON is a valid YAML subset).

Additional (shared/generic) secrets are loaded in order BEFORE the app-specific secret; later secrets override earlier ones. Missing additional secrets are silently skipped. If neither any additional nor the app-specific secret exists, the loader is skipped.

Both project id and additional secrets can also be injected programmatically. The matching env var, when set, takes priority over the option.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loadergcpsecret.New(
        loadergcpsecret.WithProjectID("my-gcp-project"),
        loadergcpsecret.WithAdditionalSecrets("generic", "keycloak"),
    )),
)
```

</details>

<details><summary>#### AWS Secrets Manager</summary>

AWS Secrets Manager loader is used to load configuration from AWS Secrets Manager.  
This is not enabled by default.

Enable AWS Secrets Manager loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/loader/external/loaderawssecrets"
)
```

| Env Value            | Description                                                                                                      | Default |
| -------------------- | ---------------------------------------------------------------------------------------------------------------- | ------- |
| `AWS_REGION`         | AWS region; if neither this nor `AWS_DEFAULT_REGION` nor the option is set, the loader is skipped (env > option) | -       |
| `AWS_DEFAULT_REGION` | Fallback region when `AWS_REGION` is unset                                                                       | -       |
| `AWS_SECRETS_PREFIX` | Prefix prepended to the configuration name to build the secret id (env > option)                                 | -       |

It calls `GetSecretValue` to fetch the secret payload (`SecretString`, falling back to `SecretBinary`).  
The configuration name is used as the secret id (optionally prefixed).  
Authentication uses the default AWS credential chain (env vars, shared config/credentials, IAM role, etc.).  
Payload is decoded as YAML by default (JSON is a valid YAML subset).

Region and prefix can also be injected programmatically. The matching env var, when set, takes priority over the option.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loaderawssecrets.New(
        loaderawssecrets.WithRegion("eu-west-1"),
        loaderawssecrets.WithPrefix("my-team"),
    )),
)
```

</details>

<details><summary>#### AWS SSM Parameter Store</summary>

AWS SSM Parameter Store loader is used to load configuration from AWS Systems Manager Parameter Store.  
This is not enabled by default.

Enable AWS SSM Parameter Store loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/loader/external/loaderawsssm"
)
```

| Env Value            | Description                                                                                                      | Default |
| -------------------- | ---------------------------------------------------------------------------------------------------------------- | ------- |
| `AWS_REGION`         | AWS region; if neither this nor `AWS_DEFAULT_REGION` nor the option is set, the loader is skipped (env > option) | -       |
| `AWS_DEFAULT_REGION` | Fallback region when `AWS_REGION` is unset                                                                       | -       |
| `AWS_SSM_PREFIX`     | Prefix prepended to the configuration name to build the parameter name (env > option)                            | -       |

It calls `GetParameter` with `WithDecryption=true` to fetch the parameter value.  
The configuration name is used as the parameter name; when a prefix is set the name becomes `/{prefix}/{name}`.  
Authentication uses the default AWS credential chain (env vars, shared config/credentials, IAM role, etc.).  
Payload is decoded as YAML by default (JSON is a valid YAML subset).

Region and prefix can also be injected programmatically. The matching env var, when set, takes priority over the option.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loaderawsssm.New(
        loaderawsssm.WithRegion("eu-west-1"),
        loaderawsssm.WithPrefix("my-team"),
    )),
)
```

</details>

<details><summary>#### Azure Key Vault</summary>

Azure Key Vault loader is used to load configuration from Azure Key Vault secrets.  
This is not enabled by default.

Enable Azure Key Vault loader importing the package.

```go
import (
    _ "github.com/rakunlabs/chu/loader/external/loaderazurekeyvault"
)
```

| Env Value             | Description                                                                                                                                     | Default |
| --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `AZURE_KEYVAULT_URL`  | Key Vault URL, e.g. `https://my-vault.vault.azure.net/`; if neither this nor `AZURE_KEYVAULT_NAME` nor the option is set, the loader is skipped | -       |
| `AZURE_KEYVAULT_NAME` | Key Vault name used to build the URL when `AZURE_KEYVAULT_URL` is unset                                                                         | -       |

It calls `GetSecret` (latest version) to fetch the secret payload.  
The configuration name is used as the secret name; `/` characters in the name are replaced with `-` because Key Vault secret names allow only alphanumeric characters and dashes.  
Authentication uses `DefaultAzureCredential` (env vars, managed identity, Azure CLI, etc.).  
Payload is decoded as YAML by default (JSON is a valid YAML subset).

The vault URL can also be injected programmatically. The `AZURE_KEYVAULT_URL` / `AZURE_KEYVAULT_NAME` env vars, when set, take priority over the option.

```go
err := chu.Load(ctx, "my-app", &cfg,
    chu.WithLoaderOption(loaderazurekeyvault.New(
        loaderazurekeyvault.WithVaultURL("https://my-vault.vault.azure.net/"),
    )),
)
```

</details>
