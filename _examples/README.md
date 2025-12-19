# Chu Examples

## Basic Example

This example for using `default`, `file` and `env` loader.

## Mix Example

For `Vault` and `Consul` loader testing.

Prepare Consul

```sh
# start consul
docker run -it --rm -p 8500:8500 --name=consul consul:latest
```

Add some data to consul

```sh
# Open http://localhost:8500
curl -X PUT http://localhost:8500/v1/kv/config/app/mix?raw -d '{"test":"99"}'
```

Prepare Vault

```sh
# start vault
docker run -it --rm --cap-add=IPC_LOCK --name=vault -p 8200:8200 -v ${PWD}/mix:/mix vault:latest
```

Get in vault container `docker exec -it vault /bin/sh` and run:

```sh
# export address for http
export VAULT_ADDR="http://127.0.0.1:8200"
# login with root token (appears in docker output)
vault login <token>
# unseal it
vault operator unseal <unsealkey>
# create kv secret engine
vault secrets enable -path=config -version=2 kv
# create policy to read
{
cat <<EOF
path "config/*" {
  capabilities = ["read", "list"]
}
EOF
} | vault policy write config-read -

# create a approle with policy and enable connection without secret_id
vault auth enable approle
vault write auth/approle/role/config-role bind_secret_id=false secret_id_bound_cidrs="127.0.0.0/8,172.17.0.0/16" policies="default","config-read"

# fill some data
vault kv put config/app/mix @mix/testdata/mix.json

# learn role-id
echo $(vault read -field=role_id auth/approle/role/config-role/role-id)
```

Run example

```sh
export VAULT_ADDR="http://127.0.0.1:8200"
export VAULT_ROLE_ID="<role-id>"
export VAULT_SECRET_BASE_PATH="config"

export CONSUL_HTTP_ADDR="localhost:8500"
export CONSUL_CONFIG_PATH_PREFIX="config"
```

```sh
go run main.go --number 2
```
