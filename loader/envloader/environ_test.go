package envloader

import (
	"testing"
)

func Test_regexIndex(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "basic",
			args: args{
				key: "_10_key",
			},
			want: 10,
		},
		{
			name: "invalid",
			args: args{
				key: "1_key",
			},
			want: -1,
		},
		{
			name: "multi",
			args: args{
				key: "_55_1_key",
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := regexIndex(tt.args.key); got != tt.want {
				t.Errorf("regexIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envHolder_MaxValue(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		e    envHolder
		args args
		want int
	}{
		{
			name: "basic",
			e: envHolder{
				"KEY_0_":  "key1",
				"KEY_1_":  "key2",
				"KEY_55_": "key55",
			},
			args: args{
				key: "KEY",
			},
			want: 55,
		},
		{
			name: "multi",
			e: envHolder{
				"DELTA_0_2_3": "key1",
				"DELTA_0_2_1": "key2",
				"DELTA_0_2_2": "key55",
			},
			args: args{
				key: "DELTA_0_2",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.MaxValue(tt.args.key); got != tt.want {
				t.Errorf("envHolder.MaxValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
