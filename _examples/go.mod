module github.com/rakunlabs/chu/example

go 1.25.3

replace github.com/rakunlabs/chu => ../

replace github.com/rakunlabs/chu/loader/external/loaderconsul => ../loader/external/loaderconsul

replace github.com/rakunlabs/chu/loader/external/loadervault => ../loader/external/loadervault

require (
	github.com/rakunlabs/chu v0.3.0
	github.com/rakunlabs/chu/loader/external/loaderconsul v0.0.0-00010101000000-000000000000
	github.com/rakunlabs/chu/loader/external/loadervault v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/hashicorp/consul/api v1.33.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-envparse v0.1.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.2.0 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.7 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hashicorp/vault/api v1.22.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/rakunlabs/logi v0.4.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/twmb/tlscfg v1.2.1 // indirect
	github.com/worldline-go/klient v0.9.13 // indirect
	github.com/worldline-go/logz v0.5.1 // indirect
	github.com/worldline-go/struct2 v1.3.1 // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20250808145144-a408d31f581a // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/time v0.12.0 // indirect
)
