package loaderawssecrets

import "io"

type Option func(*opt)

type opt struct {
	Decode func(r io.Reader, to any) error
	Prefix string
	Region string
}

func (o *opt) apply(opts ...Option) {
	for _, op := range opts {
		op(o)
	}
}

// WithDecode sets the decoder for the aws secrets loader.
//   - default is yaml decoder
func WithDecode(d func(r io.Reader, to any) error) Option {
	return func(o *opt) {
		o.Decode = d
	}
}

// WithPrefix sets a prefix prepended to the configuration name to build the
// secret id.
//   - AWS_SECRETS_PREFIX env var takes priority over this option.
func WithPrefix(prefix string) Option {
	return func(o *opt) {
		o.Prefix = prefix
	}
}

// WithRegion sets the AWS region for the secrets manager client.
//   - AWS_REGION / AWS_DEFAULT_REGION env vars take priority over this option.
func WithRegion(region string) Option {
	return func(o *opt) {
		o.Region = region
	}
}
