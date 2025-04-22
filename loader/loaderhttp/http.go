package loaderhttp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
	"github.com/worldline-go/klient"
)

type Loader struct {
	client *klient.Client
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

func (l *Loader) load(ctx context.Context, name string) ([]byte, string, error) {
	getURL, ok := loader.GetExistEnv("CONFIG_HTTP_ADDR")
	if getURL == "" || !ok {
		return nil, "", fmt.Errorf("CONFIG_HTTP_ADDR is required: %w", loader.ErrSkipLoader)
	}

	var prefix string
	if v, ok := loader.GetExistEnv("CONFIG_HTTP_PREFIX"); ok {
		prefix = strings.Trim(v, "/") + "/"
	}

	getURL = strings.TrimSuffix(getURL, "/") + "/" + prefix + strings.TrimPrefix(name, "/")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, getURL, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	var body []byte
	var contentType string
	if err := l.client.Do(req, func(r *http.Response) error {
		switch r.StatusCode {
		case http.StatusOK:
			var err error
			body, err = io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("failed to read response body: %w", err)
			}

			contentType = r.Header.Get("Content-Type")

			return nil
		case http.StatusNotFound, http.StatusNoContent:
			return fmt.Errorf("file not found: %w", loader.ErrSkipLoader)
		default:
			return klient.ErrResponse(r)
		}
	}); err != nil {
		return nil, "", fmt.Errorf("failed to do request: %w", err)
	}

	return body, contentType, nil
}

func (l *Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	vRaw, contentType, err := l.load(ctx, opt.Name)
	if err != nil {
		return err
	}

	opt.Logger.Info("config load http", "key", opt.Name)

	var mapping any

	decode := GetDecoder(contentType)

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
