package consulloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/hashicorp/go-hclog"
	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoderfile"
	"github.com/rakunlabs/chu/utils/decodermap"
)

type Loader struct {
	client *api.Client
	kv     *api.KV

	QueryOptions api.QueryOptions
	WriteOptions api.WriteOptions

	m sync.RWMutex

	// Decode for consul file to interface{}
	//  - default is yaml decoder
	Decode func(r io.Reader, to interface{}) error
}

func New(opts ...Option) *Loader {
	opt := option{}
	opt.apply(opts...)

	return &Loader{
		Decode: opt.Decode,
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

	return pair.Value, nil
}

func (l *Loader) Set(ctx context.Context, key string, value []byte) error {
	if err := l.setClient(); err != nil {
		return err
	}

	// Set the key
	pair := &api.KVPair{Key: key, Value: value}

	_, err := l.kv.Put(pair, l.WriteOptions.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	return nil
}

func (l *Loader) Delete(ctx context.Context, key string) error {
	if err := l.setClient(); err != nil {
		return err
	}

	// Delete the key
	_, err := l.kv.Delete(key, l.WriteOptions.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

// DynamicValue return a channel for getting latest value of key.
// This function will start a goroutine for watching key.
// The caller should call stop function when it is no longer needed.
func (l *Loader) DynamicValue(ctx context.Context, wg *sync.WaitGroup, key string) (<-chan []byte, func(), error) {
	if err := l.setClient(); err != nil {
		return nil, nil, err
	}

	plan, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  key,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed watch parse %w", err)
	}

	// not add any buffer, this is useful for getting latest change only
	vChannel := make(chan []byte)

	plan.HybridHandler = func(_ watch.BlockingParamVal, raw interface{}) {
		if raw == nil {
			return
		}

		v, ok := raw.(*api.KVPair)
		if ok {
			vChannel <- v.Value

			return
		}
	}

	runCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// this select-case for listen ctx done and plan run result same time
		select {
		case <-ctx.Done():
			plan.Stop()
		case <-runCh:
		}

		close(vChannel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runCh <- plan.RunWithClientAndHclog(l.client, hclog.NewNullLogger())
	}()

	return vChannel, plan.Stop, nil
}

func (l *Loader) LoadChu(ctx context.Context, to any, opts ...loader.Option) error {
	if _, ok := loader.GetExistEnv("CONSUL_HTTP_ADDR"); !ok {
		return fmt.Errorf("CONSUL_HTTP_ADDR is required: %w", loader.ErrSkipLoader)
	}

	opt := loader.NewOption(opts...)

	if err := l.setClient(); err != nil {
		return err
	}

	vRaw, err := l.Load(ctx, opt.Name)
	if err != nil {
		return err
	}

	var mapping interface{}

	decode := l.Decode
	if decode == nil {
		decode = decoderfile.Yaml{}.Decode
	}

	if err := decode(bytes.NewReader(vRaw), &mapping); err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	mapDecoder := opt.MapDecoder

	if mapDecoder == nil {
		mapDecoder = decodermap.New(
			decodermap.WithTag(opt.Tag),
			decodermap.WithHooks(opt.Hooks...),
		).Decode
	}

	if err := mapDecoder(mapping, to); err != nil {
		return fmt.Errorf("failed to map decode: %w", err)
	}

	return nil
}
