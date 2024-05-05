package file

import (
	"context"
	"reflect"
	"testing"

	"github.com/rakunlabs/chu/loader"
)

func TestLoader_Load(t *testing.T) {
	type fields struct {
		options []Option
	}
	type args struct {
		ctx  context.Context
		to   any
		opts []loader.Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				options: []Option{},
			},
			args: args{
				ctx: context.Background(),
				to: &struct {
					Host string `cfg:"host"`
				}{},
				opts: []loader.Option{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New()

			t.Setenv("CONFIG_PATH", "testdata/config.yaml")

			if err := l.LoadChu(tt.args.ctx, tt.args.to, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Loader.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("Loader.Load() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
