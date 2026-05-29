package loadergcpparameter

import (
	"io"

	"google.golang.org/api/option"
)

type Option func(*opt)

type opt struct {
	Decode       func(r io.Reader, to any) error
	ClientOption []option.ClientOption
	ProjectID    string

	// VersionFallback enables falling back to the most recently created
	// enabled version when "latest" points to a disabled version.
	VersionFallback bool
}

func (o *opt) apply(opts ...Option) {
	for _, op := range opts {
		op(o)
	}
}

// WithDecode sets the decoder for the gcp parameter loader.
//   - default is yaml decoder
func WithDecode(d func(r io.Reader, to any) error) Option {
	return func(o *opt) {
		o.Decode = d
	}
}

// WithClientOption sets the client options for the gcp parameter client.
func WithClientOption(opts ...option.ClientOption) Option {
	return func(o *opt) {
		o.ClientOption = opts
	}
}

// WithProjectID sets the GCP project id for the parameter loader.
//   - GCP_PROJECT_ID env var takes priority over this option.
func WithProjectID(projectID string) Option {
	return func(o *opt) {
		o.ProjectID = projectID
	}
}

// WithVersionFallback enables version fallback for the parameter loader.
//
// When enabled and the resolved version is "latest" (i.e. GCP_PARAMETER_VERSION
// is unset), a disabled latest version no longer fails the load: the loader
// lists the parameter versions and renders the most recently created enabled
// version instead.
//   - GCP_PARAMETER_FALLBACK env var takes priority over this option.
func WithVersionFallback(enable bool) Option {
	return func(o *opt) {
		o.VersionFallback = enable
	}
}
