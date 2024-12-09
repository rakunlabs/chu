package defaultloader

import (
	"context"
	"reflect"
	"testing"

	"github.com/rakunlabs/chu/loader"
)

func TestLoader_Load(t *testing.T) {
	type args struct {
		to   any
		opts []loader.Option
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
					Str    string `default:"hello"`
					Uint32 uint32 `default:"1"`
					Bool   bool   `default:"true"`
				}{},
			},
			want: &struct {
				Str    string `default:"hello"`
				Uint32 uint32 `default:"1"`
				Bool   bool   `default:"true"`
			}{
				Str:    "hello",
				Uint32: 1,
				Bool:   true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New()
			if err := l.LoadChu(context.Background(), tt.args.to, tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("Loader.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("Loader.Load() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
