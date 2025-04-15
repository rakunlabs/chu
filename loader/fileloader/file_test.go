package fileloader

import (
	"context"
	"reflect"
	"testing"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/utils/decodermap"
)

func TestLoader_Load(t *testing.T) {
	type args struct {
		ctx  context.Context
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
			name: "basic",
			args: args{
				ctx: context.Background(),
				to: &struct {
					Host string `cfg:"host"`
				}{},
				opts: []loader.OptionFunc{
					loader.WithName("config"),
				},
			},
			want: &struct {
				Host string `cfg:"host"`
			}{
				Host: "localhost",
			},
			wantErr: false,
		},
	}

	mapDecoder := decodermap.New().Decode

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New()()

			t.Setenv("CONFIG_PATH", "testdata/config.yaml")

			tt.args.opts = append([]loader.OptionFunc{
				loader.WithMapDecoder(mapDecoder),
			}, tt.args.opts...)

			opt := loader.NewOption(tt.args.opts...)

			if err := l.LoadChu(tt.args.ctx, tt.args.to, opt); (err != nil) != tt.wantErr {
				t.Errorf("Loader.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("Loader.Load() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
