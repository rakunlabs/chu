package loaderawsssm

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

// Loader loads configuration from AWS Systems Manager (SSM) Parameter Store.
//
// It fetches a single parameter value (with decryption) and decodes it
// (YAML by default).
type Loader struct {
	client *ssm.Client

	m sync.RWMutex

	// Prefix is prepended to the configuration name to build the parameter name.
	//  - AWS_SSM_PREFIX env var takes priority over this field.
	Prefix string

	// Region for the AWS client.
	//  - AWS_REGION / AWS_DEFAULT_REGION env vars take priority over this field.
	Region string

	// Decode for parameter payload to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...Option) loader.Loader {
	o := opt{}
	o.apply(opts...)

	return &Loader{
		Decode: o.Decode,
		Prefix: o.Prefix,
		Region: o.Region,
	}
}

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameAWSSSM
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderAWSSSM
}

func (l *Loader) SetClient(c *ssm.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *ssm.Client {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client
}

func (l *Loader) exist() bool {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client != nil
}

func (l *Loader) setClient(ctx context.Context, region string) error {
	if l.exist() {
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	var optFns []func(*config.LoadOptions) error
	if region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return fmt.Errorf("failed to load aws config: %w", err)
	}

	l.client = ssm.NewFromConfig(cfg)

	return nil
}

// LoadKey fetches a parameter value by its name (with decryption) and returns
// the raw payload bytes.
func (l *Loader) LoadKey(ctx context.Context, name string) ([]byte, error) {
	resp, err := l.client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(name),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		var nf *types.ParameterNotFound
		if errors.As(err, &nf) {
			return nil, fmt.Errorf("parameter not found [%s]: %w", name, loader.ErrSkipLoader)
		}

		return nil, fmt.Errorf("failed to get parameter: %w", err)
	}

	if resp.Parameter == nil || resp.Parameter.Value == nil {
		return nil, fmt.Errorf("parameter empty [%s]: %w", name, loader.ErrSkipLoader)
	}

	return []byte(*resp.Parameter.Value), nil
}

func (l *Loader) Load(ctx context.Context, to any, optLoader *loader.Option) error {
	// config.LoadDefaultConfig auto-reads AWS_REGION/AWS_DEFAULT_REGION; require
	// a region so the loader is skipped when AWS is not configured.
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if region == "" {
		region = l.Region
	}
	if region == "" {
		return fmt.Errorf("AWS_REGION is required: %w", loader.ErrSkipLoader)
	}

	// Env takes priority over option-set prefix.
	prefix := os.Getenv("AWS_SSM_PREFIX")
	if prefix == "" {
		prefix = l.Prefix
	}

	if err := l.setClient(ctx, region); err != nil {
		return err
	}

	name := strings.Trim(optLoader.Name, "/")
	if name == "" {
		return fmt.Errorf("parameter name is empty: %w", loader.ErrSkipLoader)
	}
	if prefix != "" {
		name = "/" + strings.Trim(prefix, "/") + "/" + name
	}

	vRaw, err := l.LoadKey(ctx, name)
	if err != nil {
		return err
	}

	optLoader.Logger.Info("config load aws ssm parameter", "name", name)

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
		).Decode
	}

	if err := mapDecoder(mapping, to); err != nil {
		return fmt.Errorf("failed to map decode: %w", err)
	}

	return nil
}
