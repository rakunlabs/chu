module github.com/rakunlabs/chu/example

go 1.27

replace github.com/rakunlabs/chu => ../

replace github.com/rakunlabs/chu/loader/external/loaderconsul => ../loader/external/loaderconsul

replace github.com/rakunlabs/chu/loader/external/loadervault => ../loader/external/loadervault

require (
	github.com/rakunlabs/chu v0.5.1
	github.com/rakunlabs/chu/loader/external/loaderconsul v0.0.0-00010101000000-000000000000
	github.com/rakunlabs/chu/loader/external/loadervault v0.0.0-00010101000000-000000000000
)

require (
	github.com/BurntSushi/toml v1.6.0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/go-jose/go-jose/v4 v4.1.6 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/hashicorp/consul/api v1.35.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-envparse v0.1.0 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-metrics v0.7.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-rootcerts v1.1.0 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.2.0 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.7 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/hashicorp/serf v0.11.0 // indirect
	github.com/hashicorp/vault/api v1.23.0 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/rakunlabs/gofret v0.2.1 // indirect
	github.com/rakunlabs/logi v0.4.6 // indirect
	github.com/rakunlabs/ok v0.1.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/twmb/tlscfg v1.3.0 // indirect
	github.com/xhit/go-str2duration/v2 v2.2.0 // indirect
	golang.org/x/exp v0.0.0-20261014143641-f0251eb4aa3e // indirect
	golang.org/x/net v0.62.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
	golang.org/x/time v0.16.0 // indirect
)
