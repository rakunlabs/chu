package loaderenv

import (
	"reflect"
	"testing"
	"time"

	"github.com/rakunlabs/chu/loader"
)

// Types for embedded struct tests - defined at package level to ensure type identity
type TestBaseConfig struct {
	Host string `cfg:"host"`
	Port int    `cfg:"port"`
}

type TestConfigWithEmbed struct {
	TestBaseConfig
	Name string `cfg:"name"`
}

type TestDeepConfig struct {
	Debug bool `cfg:"debug"`
}

type TestBaseConfigNested struct {
	TestDeepConfig
	Host string `cfg:"host"`
}

type TestConfigNestedEmbed struct {
	TestBaseConfigNested
	Name string `cfg:"name"`
}

type TestBaseConfigSimple struct {
	Host string `cfg:"host"`
}

type TestConfigWithExplicitTag struct {
	TestBaseConfigSimple `cfg:"base"`
	Name                 string `cfg:"name"`
}

type TestConfigWithEmbedPointer struct {
	*TestBaseConfigSimple
	Name string `cfg:"name"`
}

// Types for squash tag option tests
type TestConfigWithSquash struct {
	TestBaseConfig `cfg:",squash"`
	Name           string `cfg:"name"`
}

type TestConfigWithSquashPointer struct {
	*TestBaseConfigSimple `cfg:",squash"`
	Name                  string `cfg:"name"`
}

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
			name: "slice string",
			args: args{
				value: &struct {
					Test []string `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST": "test,abc",
			},
			want: &struct {
				Test []string `cfg:"test"`
			}{
				Test: []string{"test", "abc"},
			},
			wantErr: false,
		},
		{
			name: "slice string with _",
			args: args{
				value: &struct {
					Test []string `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST_0": "test,123",
			},
			want: &struct {
				Test []string `cfg:"test"`
			}{
				Test: []string{"test,123"},
			},
			wantErr: false,
		},
		{
			name: "slice numbers",
			args: args{
				value: &struct {
					Test []int `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST": "4342,123",
			},
			want: &struct {
				Test []int `cfg:"test"`
			}{
				Test: []int{4342, 123},
			},
			wantErr: false,
		},
		{
			name: "map test",
			args: args{
				value: &struct {
					Test map[string]string `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST_A": "test-1",
				"TEST_B": "test-2",
			},
			want: &struct {
				Test map[string]string `cfg:"test"`
			}{
				Test: map[string]string{
					"A": "test-1",
					"B": "test-2",
				},
			},
			wantErr: false,
		},
		{
			name: "map struct test",
			args: args{
				value: &struct {
					Test map[string]struct {
						Value string `cfg:"value"`
					} `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST_A_VALUE": "test-1",
				"TEST_B_VALUE": "test-2",
			},
			want: &struct {
				Test map[string]struct {
					Value string `cfg:"value"`
				} `cfg:"test"`
			}{
				Test: map[string]struct {
					Value string `cfg:"value"`
				}{
					"A": {Value: "test-1"},
					"B": {Value: "test-2"},
				},
			},
			wantErr: false,
		},
		{
			name: "map int struct test",
			args: args{
				value: &struct {
					Test map[int]struct {
						Value string `cfg:"value"`
					} `cfg:"test"`
				}{},
			},
			env: map[string]string{
				"TEST_1_VALUE": "test-1",
				"TEST_2_VALUE": "test-2",
			},
			want: &struct {
				Test map[int]struct {
					Value string `cfg:"value"`
				} `cfg:"test"`
			}{
				Test: map[int]struct {
					Value string `cfg:"value"`
				}{
					1: {Value: "test-1"},
					2: {Value: "test-2"},
				},
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
		{
			name: "embedded struct",
			args: args{
				value: &TestConfigWithEmbed{},
			},
			env: map[string]string{
				"HOST": "localhost",
				"PORT": "8080",
				"NAME": "test-app",
			},
			want: &TestConfigWithEmbed{
				TestBaseConfig: TestBaseConfig{
					Host: "localhost",
					Port: 8080,
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "nested embedded struct",
			args: args{
				value: &TestConfigNestedEmbed{},
			},
			env: map[string]string{
				"DEBUG": "true",
				"HOST":  "localhost",
				"NAME":  "test-app",
			},
			want: &TestConfigNestedEmbed{
				TestBaseConfigNested: TestBaseConfigNested{
					TestDeepConfig: TestDeepConfig{
						Debug: true,
					},
					Host: "localhost",
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "embedded struct with explicit tag",
			args: args{
				value: &TestConfigWithExplicitTag{},
			},
			env: map[string]string{
				"BASE_HOST": "localhost",
				"NAME":      "test-app",
			},
			want: &TestConfigWithExplicitTag{
				TestBaseConfigSimple: TestBaseConfigSimple{
					Host: "localhost",
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "embedded pointer struct",
			args: args{
				value: &TestConfigWithEmbedPointer{},
			},
			env: map[string]string{
				"HOST": "localhost",
				"NAME": "test-app",
			},
			want: &TestConfigWithEmbedPointer{
				TestBaseConfigSimple: &TestBaseConfigSimple{
					Host: "localhost",
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "squash tag option",
			args: args{
				value: &TestConfigWithSquash{},
			},
			env: map[string]string{
				"HOST": "localhost",
				"PORT": "8080",
				"NAME": "test-app",
			},
			want: &TestConfigWithSquash{
				TestBaseConfig: TestBaseConfig{
					Host: "localhost",
					Port: 8080,
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "squash tag option with pointer",
			args: args{
				value: &TestConfigWithSquashPointer{},
			},
			env: map[string]string{
				"HOST": "localhost",
				"NAME": "test-app",
			},
			want: &TestConfigWithSquashPointer{
				TestBaseConfigSimple: &TestBaseConfigSimple{
					Host: "localhost",
				},
				Name: "test-app",
			},
			wantErr: false,
		},
		{
			name: "noprefix fallback to unprefixed env",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:"log_level,noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"APP_HOST":  "localhost",
				"LOG_LEVEL": "debug",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:"log_level,noprefix"`
			}{
				Host:     "localhost",
				LogLevel: "debug",
			},
			wantErr: false,
		},
		{
			name: "noprefix prefixed takes priority",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:"log_level,noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"APP_HOST":      "localhost",
				"APP_LOG_LEVEL": "info",
				"LOG_LEVEL":     "debug",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:"log_level,noprefix"`
			}{
				Host:     "localhost",
				LogLevel: "info",
			},
			wantErr: false,
		},
		{
			name: "noprefix without prefix set works normally",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:"log_level,noprefix"`
				}{},
			},
			env: map[string]string{
				"HOST":      "localhost",
				"LOG_LEVEL": "debug",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:"log_level,noprefix"`
			}{
				Host:     "localhost",
				LogLevel: "debug",
			},
			wantErr: false,
		},
		{
			name: "noprefix field not found in either",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:"log_level,noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"APP_HOST": "localhost",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:"log_level,noprefix"`
			}{
				Host: "localhost",
			},
			wantErr: false,
		},
		{
			name: "normal field not affected by other noprefix fields",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:"log_level,noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"HOST":      "should-not-load",
				"LOG_LEVEL": "debug",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:"log_level,noprefix"`
			}{
				Host:     "",
				LogLevel: "debug",
			},
			wantErr: false,
		},
		{
			name: "noprefix with empty name uses field name",
			args: args{
				value: &struct {
					Host     string `cfg:"host"`
					LogLevel string `cfg:",noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"APP_HOST": "localhost",
				"LOGLEVEL": "debug",
			},
			want: &struct {
				Host     string `cfg:"host"`
				LogLevel string `cfg:",noprefix"`
			}{
				Host:     "localhost",
				LogLevel: "debug",
			},
			wantErr: false,
		},
		{
			name: "noprefix with nested struct",
			args: args{
				value: &struct {
					Host string `cfg:"host"`
					DB   struct {
						User     string `cfg:"user"`
						Password string `cfg:"password,noprefix"`
					} `cfg:"db,noprefix"`
				}{},
				opts: []Option{WithPrefix("APP_")},
			},
			env: map[string]string{
				"APP_HOST":    "localhost",
				"DB_USER":     "root",
				"DB_PASSWORD": "secret",
			},
			want: &struct {
				Host string `cfg:"host"`
				DB   struct {
					User     string `cfg:"user"`
					Password string `cfg:"password,noprefix"`
				} `cfg:"db,noprefix"`
			}{
				Host: "localhost",
				DB: struct {
					User     string `cfg:"user"`
					Password string `cfg:"password,noprefix"`
				}{
					User:     "root",
					Password: "secret",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			e := New(tt.args.opts...)
			if err := e.Load(t.Context(), tt.args.value, loader.NewOption()); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.value, tt.want) {
				t.Errorf("Load() = %#v, want %#v", tt.args.value, tt.want)
			}
		})
	}
}

func Test_hasTagOption(t *testing.T) {
	tests := []struct {
		name     string
		tagValue string
		option   string
		want     bool
	}{
		{
			name:     "squash option present",
			tagValue: ",squash",
			option:   "squash",
			want:     true,
		},
		{
			name:     "squash option with name",
			tagValue: "fieldname,squash",
			option:   "squash",
			want:     true,
		},
		{
			name:     "squash option with multiple options",
			tagValue: "fieldname,omitempty,squash",
			option:   "squash",
			want:     true,
		},
		{
			name:     "no squash option",
			tagValue: "fieldname",
			option:   "squash",
			want:     false,
		},
		{
			name:     "empty tag",
			tagValue: "",
			option:   "squash",
			want:     false,
		},
		{
			name:     "different option",
			tagValue: ",omitempty",
			option:   "squash",
			want:     false,
		},
		{
			name:     "noprefix option present",
			tagValue: "host,noprefix",
			option:   "noprefix",
			want:     true,
		},
		{
			name:     "noprefix option without name",
			tagValue: ",noprefix",
			option:   "noprefix",
			want:     true,
		},
		{
			name:     "noprefix not present",
			tagValue: "host",
			option:   "noprefix",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasTagOption(tt.tagValue, tt.option); got != tt.want {
				t.Errorf("hasTagOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
