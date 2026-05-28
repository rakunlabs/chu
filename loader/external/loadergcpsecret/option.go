package loadergcpsecret

import (
	"io"

	"google.golang.org/api/option"
)

type Option func(*opt)

type opt struct {
	Decode            func(r io.Reader, to any) error
	ClientOption      []option.ClientOption
	ProjectID         string
	AdditionalSecrets []string
}

func (o *opt) apply(opts ...Option) {
	for _, op := range opts {
		op(o)
	}
}

// WithDecode sets the decoder for the gcp secret loader.
//   - default is yaml decoder
func WithDecode(d func(r io.Reader, to any) error) Option {
	return func(o *opt) {
		o.Decode = d
	}
}

// WithClientOption sets the client options for the gcp secret client.
func WithClientOption(opts ...option.ClientOption) Option {
	return func(o *opt) {
		o.ClientOption = opts
	}
}

// WithProjectID sets the GCP project id for the secret loader.
//   - GCP_PROJECT_ID env var takes priority over this option.
func WithProjectID(projectID string) Option {
	return func(o *opt) {
		o.ProjectID = projectID
	}
}

// WithAdditionalSecrets sets a list of shared/generic secret ids that are
// loaded BEFORE the app-specific secret. Later secrets override earlier ones.
//
// Useful for mounting common credentials shared across multiple apps
// (e.g. "generic", "keycloak").
//
// Missing secrets are silently skipped.
//
//   - GCP_SECRET_ADDITIONAL env var (comma-separated) takes priority over this
//     option.
func WithAdditionalSecrets(names ...string) Option {
	return func(o *opt) {
		o.AdditionalSecrets = names
	}
}
