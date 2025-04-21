package consulloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

type Loader struct {
	client *api.Client
	kv     *api.KV

	QueryOptions api.QueryOptions
	WriteOptions api.WriteOptions

	m sync.RWMutex

	// Decode for consul file to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...Option) func() loader.Loader {
	return func() loader.Loader {
		opt := option{}
		opt.apply(opts...)

		return &Loader{
			Decode: opt.Decode,
		}
	}
}

func (l *Loader) SetClient(c *api.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
	l.kv = c.KV()
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
		return fmt.Errorf("failed to create consul client: %w", err)
	}

	l.client = client
	l.kv = client.KV()

	return nil
}

func (l *Loader) Load(ctx context.Context, key string) ([]byte, error) {
	if err := l.setClient(); err != nil {
		return nil, err
	}

	// Get the key
	pair, _, err := l.kv.Get(key, l.QueryOptions.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	if pair == nil {
		return nil, fmt.Errorf("key not found [%s]: %w", key, loader.ErrSkipLoader)
	}

	return pair.Value, nil
}

func (l *Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	if _, ok := loader.GetExistEnv("CONSUL_HTTP_ADDR"); !ok {
		return fmt.Errorf("CONSUL_HTTP_ADDR is required: %w", loader.ErrSkipLoader)
	}

	if err := l.setClient(); err != nil {
		return err
	}

	name := opt.Name
	if prefix, _ := loader.GetExistEnv("CONSUL_CONFIG_PATH_PREFIX"); prefix != "" {
		name = path.Join(prefix, name)
	}

	vRaw, err := l.Load(ctx, name)
	if err != nil {
		return err
	}

	opt.Logger.Info("config load consul", "key", name)

	var mapping any

	decode := l.Decode
	if decode == nil {
		decode = decoder.DecodeYaml
	}

	if err := decode(bytes.NewReader(vRaw), &mapping); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	mapDecoder := opt.MapDecoder

	if mapDecoder == nil {
		mapDecoder = decoder.New(
			decoder.WithTag(opt.Tag),
			decoder.WithHooks(opt.Hooks...),
		).Decode
	}

	if err := mapDecoder(mapping, to); err != nil {
		return fmt.Errorf("failed to map decode: %w", err)
	}

	return nil
}
