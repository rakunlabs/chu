package loaderazurekeyvault

import "io"

type Option func(*opt)

type opt struct {
	Decode   func(r io.Reader, to any) error
	VaultURL string
}

func (o *opt) apply(opts ...Option) {
	for _, op := range opts {
		op(o)
	}
}

// WithDecode sets the decoder for the azure keyvault loader.
//   - default is yaml decoder
func WithDecode(d func(r io.Reader, to any) error) Option {
	return func(o *opt) {
		o.Decode = d
	}
}

// WithVaultURL sets the Key Vault URL (e.g. https://my-vault.vault.azure.net/).
//   - AZURE_KEYVAULT_URL / AZURE_KEYVAULT_NAME env vars take priority over this
//     option.
func WithVaultURL(vaultURL string) Option {
	return func(o *opt) {
		o.VaultURL = vaultURL
	}
}
