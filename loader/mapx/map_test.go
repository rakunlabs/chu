package mapx

import (
	"context"
	"reflect"
	"testing"
)

func TestLoadSetter_Load(t *testing.T) {
	type fields struct {
		value interface{}
	}
	type args struct {
		to any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "nil value",
			fields: fields{
				value: nil,
			},
			wantErr: true,
		},
		{
			name: "basic struct",
			fields: fields{
				map[string]string{
					"host-2": "localhost",
				},
			},
			args: args{
				to: &struct {
					Host2 string `cfg:"host_2"`
				}{},
			},
			want: &struct {
				Host2 string `cfg:"host_2"`
			}{
				Host2: "localhost",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New()
			if err := l.SetValue(tt.fields.value).Load(context.Background(), tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("LoadSetter.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.args.to, tt.want) {
				t.Errorf("LoadSetter.Load() = %v, want %v", tt.args.to, tt.want)
			}
		})
	}
}
