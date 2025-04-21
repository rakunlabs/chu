package httploader

import (
	"reflect"
	"testing"

	"github.com/rakunlabs/chu/utils/decoder"
)

func TestGetDecoder(t *testing.T) {
	type args struct {
		contentType string
	}
	tests := []struct {
		name string
		args args
		want Decoder
	}{
		{
			name: "application/json",
			args: args{contentType: "application/json; charset=utf-8"},
			want: decoder.DecodeJson,
		},
		{
			name: "application/x-toml",
			args: args{contentType: "application/x-toml; charset=utf-8"},
			want: decoder.DecodeToml,
		},
		{
			name: "application/x-yaml",
			args: args{contentType: "application/x-yaml; charset=utf-8"},
			want: decoder.DecodeYaml,
		},
		{
			name: "text/yaml",
			args: args{contentType: "text/yaml; charset=utf-8"},
			want: decoder.DecodeYaml,
		},
		{
			name: "",
			args: args{contentType: "text/x-yaml; charset=utf-8"},
			want: decoder.DecodeYaml,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDecoder(tt.args.contentType)

			if reflect.ValueOf(got).Pointer() != reflect.ValueOf(tt.want).Pointer() {
				t.Errorf("GetDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}
