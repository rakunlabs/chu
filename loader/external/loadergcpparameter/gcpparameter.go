package loadergcpparameter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	parametermanager "cloud.google.com/go/parametermanager/apiv1"
	"cloud.google.com/go/parametermanager/apiv1/parametermanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

// Loader loads configuration from GCP Parameter Manager.
//
// It fetches a single parameter version's rendered payload and decodes it
// (YAML by default).
type Loader struct {
	client *parametermanager.Client

	clientOption []option.ClientOption

	// ProjectID falls back to GCP_PROJECT_ID env var.
	//  - GCP_PROJECT_ID env var takes priority over this field.
	ProjectID string

	m sync.RWMutex

	// Decode for parameter payload to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...Option) loader.Loader {
	o := opt{}
	o.apply(opts...)

	return &Loader{
		Decode:       o.Decode,
		clientOption: o.ClientOption,
		ProjectID:    o.ProjectID,
	}
}

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameGCPParameter
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderGCPParameter
}

func (l *Loader) SetClient(c *parametermanager.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *parametermanager.Client {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client
}

func (l *Loader) exist() bool {
	l.m.RLock()
	defer l.m.RUnlock()

	return l.client != nil
}

func (l *Loader) setClient(ctx context.Context, location string) error {
	if l.exist() {
		return nil
	}

	l.m.Lock()
	defer l.m.Unlock()

	clientOpts := l.clientOption
	// For regional (non-global) endpoints, the client needs the regional endpoint.
	if location != "" && location != "global" {
		clientOpts = append([]option.ClientOption{
			option.WithEndpoint(fmt.Sprintf("parametermanager.%s.rep.googleapis.com:443", location)),
		}, clientOpts...)
	}

	client, err := parametermanager.NewClient(ctx, clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create gcp parameter manager client: %w", err)
	}

	l.client = client

	return nil
}

// LoadKey loads a parameter version by its full resource name and returns the
// rendered payload bytes.
//
// resourceName must be in the form:
//
//	projects/{project}/locations/{location}/parameters/{parameter}/versions/{version}
func (l *Loader) LoadKey(ctx context.Context, resourceName string) ([]byte, error) {
	req := &parametermanagerpb.RenderParameterVersionRequest{
		Name: resourceName,
	}

	resp, err := l.client.RenderParameterVersion(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("parameter not found [%s]: %w", resourceName, loader.ErrSkipLoader)
		}

		return nil, fmt.Errorf("failed to render parameter version: %w", err)
	}

	return resp.GetRenderedPayload(), nil
}

func (l *Loader) Load(ctx context.Context, to any, optLoader *loader.Option) error {
	// Env takes priority over option-set ProjectID.
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = l.ProjectID
	}
	if projectID == "" {
		return fmt.Errorf("GCP_PROJECT_ID is required: %w", loader.ErrSkipLoader)
	}

	location := os.Getenv("GCP_PARAMETER_LOCATION")
	if location == "" {
		location = "global"
	}

	version := os.Getenv("GCP_PARAMETER_VERSION")
	if version == "" {
		version = "latest"
	}

	if err := l.setClient(ctx, location); err != nil {
		return err
	}

	// "/" is not allowed in GCP parameter ids; replace with "-".
	parameterID := strings.ReplaceAll(strings.Trim(optLoader.Name, "/"), "/", "-")
	if parameterID == "" {
		return fmt.Errorf("parameter name is empty: %w", loader.ErrSkipLoader)
	}

	resourceName := fmt.Sprintf(
		"projects/%s/locations/%s/parameters/%s/versions/%s",
		projectID, location, parameterID, version,
	)

	vRaw, err := l.LoadKey(ctx, resourceName)
	if err != nil {
		return err
	}

	optLoader.Logger.Info("config load gcp parameter", "name", resourceName)

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
