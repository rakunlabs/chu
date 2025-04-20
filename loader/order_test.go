package loader

import (
	"reflect"
	"testing"
)

func Test_orderLoaders(t *testing.T) {
	type args struct {
		loaders map[string]LoadHolder
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple line",
			args: args{
				loaders: map[string]LoadHolder{
					"d": {},
					"b": {
						Order: &Order{After: []string{"a"}, Before: []string{"c"}},
					},
					"a": {
						Order: &Order{Before: []string{"b"}},
					},
					"c": {
						Order: &Order{After: []string{"b", "a"}},
					},
					"e": {
						Order: &Order{Before: []string{"c"}},
					},
					"f": {
						Order: &Order{Before: []string{"e"}},
					},
				},
			},
			want: []string{"a", "f", "b", "e", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderLoaders(tt.args.loaders); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("orderLoaders() = %v, want %v", got, tt.want)
			}
		})
	}
}
