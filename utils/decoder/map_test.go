package decoder

import (
	"reflect"
	"testing"
	"time"

	"github.com/rakunlabs/gofret"
	"github.com/xhit/go-str2duration/v2"
)

func TestMapDecode(t *testing.T) {
	tests := []struct {
		name string
		new  *Map
		in   map[string]any
		want any
	}{
		{
			name: "weak types and loose keys",
			new:  New(),
			in:   map[string]any{"MAX-RETRY": "3"},
			want: struct {
				MaxRetry int `cfg:"max_retry"`
			}{MaxRetry: 3},
		},
		{
			name: "duration supports days with conversion",
			new:  New(WithCodecOptions(gofret.Between(str2duration.ParseDuration))),
			in:   map[string]any{"timeout": "2d"},
			want: struct {
				Timeout time.Duration `cfg:"timeout"`
			}{Timeout: 48 * time.Hour},
		},
		{
			name: "foreign tag option is ignored",
			new:  New(),
			in:   map[string]any{"name": "service"},
			want: struct {
				Name string `cfg:"name,noprefix"`
			}{Name: "service"},
		},
		{
			name: "squash aliases inline",
			new:  New(),
			in:   map[string]any{"name": "service"},
			want: struct {
				Nested struct {
					Name string `cfg:"name"`
				} `cfg:",squash"`
			}{Nested: struct {
				Name string `cfg:"name"`
			}{Name: "service"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := reflect.New(reflect.TypeOf(tt.want))
			if err := tt.new.Decode(tt.in, out.Interface()); err != nil {
				t.Fatalf("Decode() error = %v", err)
			}

			if got := out.Elem().Interface(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Decode() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMapDecodeLegacyHooks(t *testing.T) {
	type config struct {
		Value string `cfg:"value"`
	}

	var calls int
	decoder := New(WithHooks(
		func(gofret.HookCtx) (any, error) {
			calls++
			return nil, gofret.ErrSkip
		},
		func(ctx gofret.HookCtx) (any, error) {
			if ctx.From.Kind() == reflect.String && ctx.To.Kind() == reflect.String {
				return "hooked-" + ctx.Data().(string), nil
			}

			return nil, gofret.ErrSkip
		},
	))

	var got config
	if err := decoder.Decode(map[string]any{"value": "ok"}, &got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if got.Value != "hooked-ok" {
		t.Fatalf("Decode() value = %q, want %q", got.Value, "hooked-ok")
	}
	if calls == 0 {
		t.Fatal("Decode() did not call hook")
	}
}
