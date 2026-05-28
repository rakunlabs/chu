package loadergcpsecret

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

// Loader loads configuration from GCP Secret Manager.
//
// For each call it loads:
//  1. Each secret in AdditionalSecrets (e.g. "generic"), silently skipping
//     missing ones.
//  2. The app-specific secret named by opt.Name.
//
// Later secrets override earlier ones.
//
// Each secret's payload is decoded as YAML by default.
type Loader struct {
	client *secretmanager.Client

	clientOption []option.ClientOption

	// ProjectID falls back to GCP_PROJECT_ID env var.
	//  - GCP_PROJECT_ID env var takes priority over this field.
	ProjectID string

	// AdditionalSecrets are loaded BEFORE the app-specific secret. Missing
	// ones are silently skipped.
	//  - GCP_SECRET_ADDITIONAL env var (comma-separated) takes priority over
	//    this field.
	AdditionalSecrets []string

	m sync.RWMutex

	// Decode for secret payload to any
	//  - default is yaml decoder
	Decode func(r io.Reader, to any) error
}

func New(opts ...Option) loader.Loader {
	o := opt{}
	o.apply(opts...)

	return &Loader{
		Decode:            o.Decode,
		clientOption:      o.ClientOption,
		ProjectID:         o.ProjectID,
		AdditionalSecrets: o.AdditionalSecrets,
	}
}

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameGCPSecret
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderGCPSecret
}

func (l *Loader) SetClient(c *secretmanager.Client) {
	l.m.Lock()
	defer l.m.Unlock()

	l.client = c
}

func (l *Loader) Client() *secretmanager.Client {
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
			option.WithEndpoint(fmt.Sprintf("secretmanager.%s.rep.googleapis.com:443", location)),
		}, clientOpts...)
	}

	client, err := secretmanager.NewClient(ctx, clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create gcp secret manager client: %w", err)
	}

	l.client = client

	return nil
}

// LoadKey loads a secret version by its full resource name and returns the
// payload bytes. NotFound errors are wrapped with loader.ErrSkipLoader.
//
// resourceName must be in one of these forms:
//
//	projects/{project}/secrets/{secret}/versions/{version}                       (global)
//	projects/{project}/locations/{location}/secrets/{secret}/versions/{version}  (regional)
func (l *Loader) LoadKey(ctx context.Context, resourceName string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceName,
	}

	resp, err := l.client.AccessSecretVersion(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, fmt.Errorf("secret not found [%s]: %w", resourceName, loader.ErrSkipLoader)
		}

		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	return resp.GetPayload().GetData(), nil
}

// buildResourceName constructs the full secret version resource name.
func buildResourceName(projectID, location, secretID, version string) string {
	if location == "" || location == "global" {
		return fmt.Sprintf(
			"projects/%s/secrets/%s/versions/%s",
			projectID, secretID, version,
		)
	}

	return fmt.Sprintf(
		"projects/%s/locations/%s/secrets/%s/versions/%s",
		projectID, location, secretID, version,
	)
}

// sanitizeSecretID converts a name like "finops/br_sr" into a GCP-safe id.
// GCP secret ids do not allow "/".
func sanitizeSecretID(name string) string {
	return strings.ReplaceAll(strings.Trim(name, "/"), "/", "-")
}

// parseAdditionalEnv parses a comma-separated list, trimming whitespace and
// dropping empty entries.
func parseAdditionalEnv(v string) []string {
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}

	return out
}

// loadAndDecode fetches one secret and decodes its payload into 'to'.
// Returns (true, nil) on success, (false, nil) on NotFound (silent skip),
// (false, err) on any other error.
func (l *Loader) loadAndDecode(ctx context.Context, resourceName string, to any, optLoader *loader.Option) (bool, error) {
	vRaw, err := l.LoadKey(ctx, resourceName)
	if err != nil {
		if errors.Is(err, loader.ErrSkipLoader) {
			return false, nil
		}

		return false, err
	}

	optLoader.Logger.Info("config load gcp secret", "name", resourceName)

	var mapping any

	decode := l.Decode
	if decode == nil {
		decode = decoder.DecodeYaml
	}

	if err := decode(bytes.NewReader(vRaw), &mapping); err != nil {
		return false, fmt.Errorf("failed to decode: %w", err)
	}

	mapDecoder := optLoader.MapDecoder

	if mapDecoder == nil {
		mapDecoder = decoder.New(
			decoder.WithTag(optLoader.Tag),
			decoder.WithHooks(optLoader.Hooks...),
		).Decode
	}

	if err := mapDecoder(mapping, to); err != nil {
		return false, fmt.Errorf("failed to map decode: %w", err)
	}

	return true, nil
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

	location := os.Getenv("GCP_SECRET_LOCATION")
	if location == "" {
		location = "global"
	}

	version := os.Getenv("GCP_SECRET_VERSION")
	if version == "" {
		version = "latest"
	}

	// Env takes priority over option-set AdditionalSecrets.
	var additional []string
	if v := os.Getenv("GCP_SECRET_ADDITIONAL"); v != "" {
		additional = parseAdditionalEnv(v)
	} else {
		additional = l.AdditionalSecrets
	}

	if err := l.setClient(ctx, location); err != nil {
		return err
	}

	loadedAny := false

	// Load additional secrets first; missing ones are silently skipped.
	for _, name := range additional {
		secretID := sanitizeSecretID(name)
		if secretID == "" {
			continue
		}

		resourceName := buildResourceName(projectID, location, secretID, version)

		loaded, err := l.loadAndDecode(ctx, resourceName, to, optLoader)
		if err != nil {
			return err
		}
		if loaded {
			loadedAny = true
		}
	}

	// Load app-specific secret.
	secretID := sanitizeSecretID(optLoader.Name)
	if secretID == "" {
		if loadedAny {
			return nil
		}

		return fmt.Errorf("secret name is empty: %w", loader.ErrSkipLoader)
	}

	resourceName := buildResourceName(projectID, location, secretID, version)

	loaded, err := l.loadAndDecode(ctx, resourceName, to, optLoader)
	if err != nil {
		return err
	}

	if !loaded && !loadedAny {
		return fmt.Errorf("secret not found [%s]: %w", resourceName, loader.ErrSkipLoader)
	}

	return nil
}
