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
	okclient "github.com/rakunlabs/ok"
)

type Loader struct {
	client *okclient.Client
}

func New(opts ...okclient.OptionClientFn) loader.Loader {
	client, err := okclient.New(opts...)
	if err != nil {
		panic(err)
	}

	return &Loader{
		client: client,
	}
}

func (l *Loader) load(ctx context.Context, name string) ([]byte, string, error) {
	getURL, ok := loader.GetExistEnv("CONFIG_HTTP_ADDR")
	if getURL == "" || !ok {
		return nil, "", fmt.Errorf("CONFIG_HTTP_ADDR is required: %w", loader.ErrSkipLoader)
	}

	var suffix string
	if v, ok := loader.GetExistEnv("CONFIG_HTTP_SUFFIX"); ok {
		suffix = "/" + strings.Trim(v, "/")
	}

	getURL = strings.TrimSuffix(getURL, "/") + "/" + strings.Trim(name, "/") + suffix

	if query, ok := loader.GetExistEnv("CONFIG_HTTP_QUERY"); ok && query != "" {
		getURL += "?" + query
	}

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
			return okclient.ErrResponse(r)
		}
	}); err != nil {
		return nil, "", fmt.Errorf("failed to do request: %w", err)
	}

	return body, contentType, nil
}

func (l *Loader) Load(ctx context.Context, to any, opt *loader.Option) error {
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

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameHTTP
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderHTTP
}
