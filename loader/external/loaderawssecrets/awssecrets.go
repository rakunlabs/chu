package loaderawssecrets

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
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

// Loader loads configuration from AWS Secrets Manager.
//
// It fetches a single secret value and decodes it (YAML by default).
type Loader struct {
	client *secretsmanager.Client

	m sync.RWMutex

	// Prefix is prepended to the configuration name to build the secret id.
	//  - AWS_SECRETS_PREFIX env var takes priority over this field.
	Prefix string

	// Region for the AWS client.
	//  - AWS_REGION / AWS_DEFAULT_REGION env vars take priority over this field.
	Region string

	// Decode for secret payload to any
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
	return loader.NameAWSSecrets
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderAWSSecrets
}

func (l *Loader) SetClient(c *secretsmanager.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *secretsmanager.Client {
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

	l.client = secretsmanager.NewFromConfig(cfg)

	return nil
}

// LoadKey fetches a secret value by its id and returns the raw payload bytes.
func (l *Loader) LoadKey(ctx context.Context, secretID string) ([]byte, error) {
	resp, err := l.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	})
	if err != nil {
		var nf *types.ResourceNotFoundException
		if errors.As(err, &nf) {
			return nil, fmt.Errorf("secret not found [%s]: %w", secretID, loader.ErrSkipLoader)
		}

		return nil, fmt.Errorf("failed to get secret value: %w", err)
	}

	if resp.SecretString != nil {
		return []byte(*resp.SecretString), nil
	}

	return resp.SecretBinary, nil
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
	prefix := os.Getenv("AWS_SECRETS_PREFIX")
	if prefix == "" {
		prefix = l.Prefix
	}

	if err := l.setClient(ctx, region); err != nil {
		return err
	}

	secretID := strings.Trim(optLoader.Name, "/")
	if secretID == "" {
		return fmt.Errorf("secret name is empty: %w", loader.ErrSkipLoader)
	}
	if prefix != "" {
		secretID = strings.Trim(prefix, "/") + "/" + secretID
	}

	vRaw, err := l.LoadKey(ctx, secretID)
	if err != nil {
		return err
	}

	optLoader.Logger.Info("config load aws secret", "name", secretID)

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
