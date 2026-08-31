package loaderazurekeyvault

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

// Loader loads configuration from Azure Key Vault.
//
// It fetches a single secret value and decodes it (YAML by default).
type Loader struct {
	client *azsecrets.Client

	m sync.RWMutex

	// VaultURL is the Key Vault URL (e.g. https://my-vault.vault.azure.net/).
	//  - AZURE_KEYVAULT_URL / AZURE_KEYVAULT_NAME env vars take priority over
	//    this field.
	VaultURL string

	// Decode for secret payload to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...Option) loader.Loader {
	o := opt{}
	o.apply(opts...)

	return &Loader{
		Decode:   o.Decode,
		VaultURL: o.VaultURL,
	}
}

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameAzureKeyVault
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderAzureKeyVault
}

func (l *Loader) SetClient(c *azsecrets.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *azsecrets.Client {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client
}

func (l *Loader) exist() bool {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client != nil
}

func (l *Loader) setClient(vaultURL string) error {
	if l.exist() {
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to create azure credential: %w", err)
	}

	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		return fmt.Errorf("failed to create azure keyvault client: %w", err)
	}

	l.client = client

	return nil
}

// LoadKey fetches the latest version of a secret by name and returns the raw
// payload bytes.
func (l *Loader) LoadKey(ctx context.Context, name string) ([]byte, error) {
	// Empty version selects the latest enabled version.
	resp, err := l.client.GetSecret(ctx, name, "", nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("secret not found [%s]: %w", name, loader.ErrSkipLoader)
		}

		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	if resp.Value == nil {
		return nil, fmt.Errorf("secret empty [%s]: %w", name, loader.ErrSkipLoader)
	}

	return []byte(*resp.Value), nil
}

func (l *Loader) Load(ctx context.Context, to any, optLoader *loader.Option) error {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	if vaultURL == "" {
		if name := os.Getenv("AZURE_KEYVAULT_NAME"); name != "" {
			vaultURL = fmt.Sprintf("https://%s.vault.azure.net/", name)
		}
	}
	if vaultURL == "" {
		vaultURL = l.VaultURL
	}
	if vaultURL == "" {
		return fmt.Errorf("AZURE_KEYVAULT_URL or AZURE_KEYVAULT_NAME is required: %w", loader.ErrSkipLoader)
	}

	if err := l.setClient(vaultURL); err != nil {
		return err
	}

	// Key Vault secret names allow only alphanumeric characters and dashes;
	// "/" is not allowed, so replace it with "-".
	secretName := strings.ReplaceAll(strings.Trim(optLoader.Name, "/"), "/", "-")
	if secretName == "" {
		return fmt.Errorf("secret name is empty: %w", loader.ErrSkipLoader)
	}

	vRaw, err := l.LoadKey(ctx, secretName)
	if err != nil {
		return err
	}

	optLoader.Logger.Info("config load azure keyvault secret", "name", secretName)

	var mapping any

	decode := l.Decode
	if decode == nil {
		decode = decoder.DecodeYaml
	}

	if err := decode(bytes.NewReader(vRaw), &mapping); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	mapDecoder := optLoader.MapDecoder

	if mapDecoder == nil {
		mapDecoder = decoder.New(
			decoder.WithTag(optLoader.Tag),
			decoder.WithHooks(optLoader.Hooks...),
			decoder.WithCodecOptions(optLoader.CodecOptions...),
		).Decode
	}

	if err := mapDecoder(mapping, to); err != nil {
		return fmt.Errorf("failed to map decode: %w", err)
	}

	return nil
}
