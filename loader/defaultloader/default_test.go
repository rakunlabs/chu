package defaultloader

import (
	"reflect"
	"testing"
	"time"

	"github.com/rakunlabs/chu/loader"
)

func TestLoader_Load(t *testing.T) {
	type args struct {
		to   any
		opts []loader.OptionFunc
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				to: &struct {
					Str     string `default:"hello"`
					Uint32  uint32 `default:"1"`
					Bool    bool   `default:"true"`
					Defined string
					Ok      time.Duration `default:"2d"`
				}{
					Defined: "defined",
				},
				opts: []loader.OptionFunc{
					loader.WithHooks(loader.HookTimeDuration),
				},
			},
			want: &struct {
				Str     string `default:"hello"`
				Uint32  uint32 `default:"1"`
				Bool    bool   `default:"true"`
				Defined string
				Ok      time.Duration `default:"2d"`
			}{
				Str:     "hello",
				Uint32:  1,
				Bool:    true,
				Defined: "defined",
				Ok:      2 * 24 * time.Hour,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New()
			opt := loader.NewOption(tt.args.opts...)
			if err := l.LoadChu(t.Context(), tt.args.to, opt); (err != nil) != tt.wantErr {
				t.Errorf("Loader.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("Loader.Load() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
