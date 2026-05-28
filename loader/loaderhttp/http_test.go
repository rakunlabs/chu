package loaderhttp

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decoder"
)

func TestLoader_LoadChu(t *testing.T) {
	type Config struct {
		Host string `cfg:"host"`
		Port int    `cfg:"port"`
	}

	// Prepare YAML config
	configYAML := "host: testhost\nport: 1234\n"

	// Start a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/config" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(configYAML))
	}))
	defer ts.Close()

	// Set CONFIG_HTTP_ADDR to the test server's URL
	t.Setenv("CONFIG_HTTP_ADDR", ts.URL)

	cfg := &Config{}
	mapDecoder := decoder.New().Decode
	opt := loader.NewOption(
		loader.WithName("config"),
		loader.WithMapDecoder(mapDecoder),
	)

	l := New()
	err := l.Load(t.Context(), cfg, opt)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{Host: "testhost", Port: 1234}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %#v, want %#v", cfg, want)
	}
}

func TestLoader_LoadChuWithQuery(t *testing.T) {
	type Config struct {
		Host string `cfg:"host"`
		Port int    `cfg:"port"`
	}

	configYAML := "host: prodhost\nport: 8080\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/config" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// Verify multiple query parameters are received
		if r.URL.Query().Get("env") != "prod" {
			http.Error(w, "expected env=prod", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("region") != "us-east-1" {
			http.Error(w, "expected region=us-east-1", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(configYAML))
	}))
	defer ts.Close()

	t.Setenv("CONFIG_HTTP_ADDR", ts.URL)
	t.Setenv("CONFIG_HTTP_QUERY", "env=prod&region=us-east-1")

	cfg := &Config{}
	mapDecoder := decoder.New().Decode
	opt := loader.NewOption(
		loader.WithName("config"),
		loader.WithMapDecoder(mapDecoder),
	)

	l := New()
	err := l.Load(t.Context(), cfg, opt)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{Host: "prodhost", Port: 8080}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %#v, want %#v", cfg, want)
	}
}

func TestLoader_LoadChuWithVersion(t *testing.T) {
	type Config struct {
		Host string `cfg:"host"`
		Port int    `cfg:"port"`
	}

	configYAML := "host: verhost\nport: 9090\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		if r.URL.Query().Get("version") != "1.2.3" {
			http.Error(w, "expected version=1.2.3", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(configYAML))
	}))
	defer ts.Close()

	t.Setenv("CONFIG_HTTP_ADDR", ts.URL)

	cfg := &Config{}
	mapDecoder := decoder.New().Decode
	opt := loader.NewOption(
		loader.WithName("config"),
		loader.WithMapDecoder(mapDecoder),
		loader.WithVersion("1.2.3"),
	)

	l := New()
	err := l.Load(t.Context(), cfg, opt)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{Host: "verhost", Port: 9090}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %#v, want %#v", cfg, want)
	}
}

// TestLoader_LoadChuQueryOverridesVersion ensures CONFIG_HTTP_QUERY env wins
// over the version option when both are set.
func TestLoader_LoadChuQueryOverridesVersion(t *testing.T) {
	type Config struct {
		Host string `cfg:"host"`
		Port int    `cfg:"port"`
	}

	configYAML := "host: qhost\nport: 7070\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/config" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		// version query should NOT be present
		if r.URL.Query().Has("version") {
			http.Error(w, "version query must not be set", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("env") != "prod" {
			http.Error(w, "expected env=prod", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(configYAML))
	}))
	defer ts.Close()

	t.Setenv("CONFIG_HTTP_ADDR", ts.URL)
	t.Setenv("CONFIG_HTTP_QUERY", "env=prod")

	cfg := &Config{}
	mapDecoder := decoder.New().Decode
	opt := loader.NewOption(
		loader.WithName("config"),
		loader.WithMapDecoder(mapDecoder),
		loader.WithVersion("1.2.3"),
	)

	l := New()
	err := l.Load(t.Context(), cfg, opt)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := &Config{Host: "qhost", Port: 7070}
	if !reflect.DeepEqual(cfg, want) {
		t.Errorf("Load() = %#v, want %#v", cfg, want)
	}
}
