package envloader

import (
	"reflect"
	"testing"
	"time"

	"github.com/rakunlabs/chu/loader"
)

func TestLoad(t *testing.T) {
	type args struct {
		value any
		opts  []Option
	}
	tests := []struct {
		name    string
		args    args
		env     map[string]string
		want    any
		wantErr bool
	}{
		{
			name: "nil value",
			args: args{
				value: nil,
			},
			wantErr: true,
		},
		{
			name: "basic struct",
			args: args{
				value: &struct {
					Host string `cfg:"host"`
				}{},
			},
			env: map[string]string{
				"HOST": "localhost",
			},
			want: &struct {
				Host string `cfg:"host"`
			}{
				Host: "localhost",
			},
			wantErr: false,
		},
		{
			name: "nested struct",
			args: args{
				value: &struct {
					Host string `cfg:"host"`
					Port int    `cfg:"port"`
					DB   struct {
						User     string `cfg:"user"`
						Password string `cfg:"password"`
					} `cfg:"db"`
					Meta []struct {
						Key   string `cfg:"key"`
						Value string `cfg:"value"`
					} `cfg:"meta"`
					Delta [][][]int
				}{},
			},
			env: map[string]string{
				"HOST":         "localhost",
				"PORT":         "3306",
				"DB_USER":      "root",
				"DB_PASSWORD":  "password",
				"META_0_KEY":   "key1",
				"META_0_VALUE": "value1",
				"META_3_KEY":   "key3",
				"META_3_VALUE": "value3",
				"DELTA_0_2_2":  "2",
				"DELTA_0_2_3":  "3",
			},
			want: &struct {
				Host string `cfg:"host"`
				Port int    `cfg:"port"`
				DB   struct {
					User     string `cfg:"user"`
					Password string `cfg:"password"`
				} `cfg:"db"`
				Meta []struct {
					Key   string `cfg:"key"`
					Value string `cfg:"value"`
				} `cfg:"meta"`
				Delta [][][]int
			}{
				Host: "localhost",
				Port: 3306,
				DB: struct {
					User     string `cfg:"user"`
					Password string `cfg:"password"`
				}{
					User:     "root",
					Password: "password",
				},
				Meta: []struct {
					Key   string `cfg:"key"`
					Value string `cfg:"value"`
				}{
					{
						Key:   "key1",
						Value: "value1",
					},
					{Key: "", Value: ""},
					{Key: "", Value: ""},
					{
						Key:   "key3",
						Value: "value3",
					},
				},
				Delta: [][][]int{
					{
						nil,
						nil,
						{0, 0, 2, 3},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nested struct with pointer",
			args: args{
				value: &struct {
					Host string `cfg:"host"`
					Port int    `cfg:"port"`
					DB   *struct {
						User     string `cfg:"user"`
						Password string `cfg:"password"`
					} `cfg:"db"`
					DBPrivate *struct {
						User     string `cfg:"user"`
						Password string `cfg:"password"`
					} `cfg:"-"`
					Test *bool
				}{},
			},
			env: map[string]string{
				"HOST":        "localhost",
				"PORT":        "3306",
				"DB_USER":     "root",
				"DB_PASSWORD": "password",
			},
			want: &struct {
				Host string `cfg:"host"`
				Port int    `cfg:"port"`
				DB   *struct {
					User     string `cfg:"user"`
					Password string `cfg:"password"`
				} `cfg:"db"`
				DBPrivate *struct {
					User     string `cfg:"user"`
					Password string `cfg:"password"`
				} `cfg:"-"`
				Test *bool
			}{
				Host: "localhost",
				Port: 3306,
				DB: &struct {
					User     string `cfg:"user"`
					Password string `cfg:"password"`
				}{
					User:     "root",
					Password: "password",
				},
				Test: nil,
			},
			wantErr: false,
		},
		{
			name: "special type",
			args: args{
				value: &struct {
					Duration time.Duration `cfg:"duration"`
				}{},
				opts: []Option{
					WithHooks(func(input, output reflect.Type, data interface{}) (interface{}, error) {
						if input.Kind() == reflect.String && output == reflect.TypeFor[time.Duration]() {
							d, err := time.ParseDuration(data.(string))
							if err != nil {
								return nil, err
							}

							return d, nil
						}

						return data, nil
					}),
				},
			},
			env: map[string]string{
				"DURATION": "1s",
			},
			want: &struct {
				Duration time.Duration `cfg:"duration"`
			}{
				Duration: time.Second,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			e := New(tt.args.opts...)()
			if err := e.LoadChu(t.Context(), tt.args.value, loader.NewOption()); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.value, tt.want) {
				t.Errorf("Load() = %#v, want %#v", tt.args.value, tt.want)
			}
		})
	}
}
