package vaultloader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/vault/api"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decodermap"
)

type Loader struct {
	client          *api.Client
	AppRoleBasePath string

	m sync.RWMutex
}

func New() func() loader.Loader {
	return func() loader.Loader {
		return &Loader{}
	}
}

func (l *Loader) SetClient(c *api.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *api.Client {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client
}

func (l *Loader) exist() bool {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client != nil
}

func (l *Loader) setClient() error {
	if l.exist() {
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	// Get a new client
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	l.client = client

	return nil
}

func (l *Loader) Login(ctx context.Context) error {
	if err := l.setClient(); err != nil {
		return err
	}

	// A combination of a Role ID and Secret ID is required to log in to Vault with an AppRole.
	// First, let's get the role ID given to us by our Vault administrator.
	roleID := os.Getenv("VAULT_ROLE_ID")
	if roleID == "" {
		return nil
	}

	// check default path
	appRoleBasePath := l.AppRoleBasePath
	if appRoleBasePath == "" {
		appRoleBasePath = os.Getenv("VAULT_APPROLE_BASE_PATH")
	}

	if appRoleBasePath == "" {
		appRoleBasePath = "auth/approle/login"
	}

	secret, err := l.client.Logical().WriteWithContext(ctx, appRoleBasePath, map[string]interface{}{
		"role_id":   roleID,
		"secret_id": os.Getenv("VAULT_ROLE_SECRET"),
	})
	if err != nil {
		return fmt.Errorf("failed to login to vault: %w", err)
	}

	// Set the token
	l.client.SetToken(secret.Auth.ClientToken)

	return nil
}

// Load loads a key from the vault.
//   - first login to vault
func (l *Loader) Load(ctx context.Context, mountPath string, key string) (map[string]interface{}, error) {
	// Get the key
	secret, err := l.client.KVv2(mountPath).Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	return secret.Data, nil
}

func (l *Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	if _, ok := loader.GetExistEnv("VAULT_ADDR", "VAULT_AGENT_ADDR"); !ok {
		return fmt.Errorf("VAULT_ADDR or VAULT_AGENT_ADDR is required: %w", loader.ErrSkipLoader)
	}

	mountPath := os.Getenv("CONFIG_VAULT_PREFIX")
	if mountPath == "" {
		return errors.New("CONFIG_VAULT_PREFIX is required as mount path")
	}

	if err := l.Login(ctx); err != nil {
		return err
	}

	v, err := l.Load(ctx, mountPath, opt.Name)
	if err != nil {
		return err
	}

	mapDecoder := opt.MapDecoder

	if mapDecoder == nil {
		mapDecoder = decodermap.New(
			decodermap.WithTag(opt.Tag),
			decodermap.WithHooks(opt.Hooks...),
		).Decode
	}

	if err := mapDecoder(v, to); err != nil {
		return fmt.Errorf("failed to map decode: %w", err)
	}

	return nil
}
