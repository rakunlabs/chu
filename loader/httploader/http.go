package httploader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoderfile"
	"github.com/rakunlabs/chu/utils/decodermap"
	"github.com/worldline-go/klient"
)

type Loader struct {
	client *klient.Client

	// Decode to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...klient.OptionClientFn) func() loader.Loader {
	return func() loader.Loader {
		opts = append([]klient.OptionClientFn{
			klient.WithDisableBaseURLCheck(true),
		}, opts...)

		client, err := klient.New(opts...)
		if err != nil {
			panic(err)
		}

		return &Loader{
			client: client,
		}
	}
}

func (l *Loader) load(ctx context.Context) ([]byte, error) {
	getURL, ok := loader.GetExistEnv("CONFIG_HTTP_ADDR")
	if !ok {
		return nil, fmt.Errorf("CONFIG_HTTP_ADDR is required: %w", loader.ErrSkipLoader)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var body []byte
	if err := l.client.Do(req, func(r *http.Response) error {
		if r.StatusCode != http.StatusOK {
			return klient.ErrResponse(r)
		}

		var err error
		body, err = io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	return body, nil
}

func (l *Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	vRaw, err := l.load(ctx)
	if err != nil {
		return err
	}

	var mapping any

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
