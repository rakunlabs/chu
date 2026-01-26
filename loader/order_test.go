package loader

import (
	"context"
	"reflect"
	"testing"
)

type testLoader struct {
	Name  LoaderName
	Order Order
}

func (d testLoader) LoadName() LoaderName {
	return d.Name
}

func (d testLoader) LoadOrder() Order {
	return d.Order
}

func (d testLoader) Load(ctx context.Context, to any, opt *Option) error {
	return nil
}

func Test_orderLoaders(t *testing.T) {
	type args struct {
		loaders map[LoaderName]Loader
	}
	tests := []struct {
		name string
		args args
		want []LoaderName
	}{
		{
			name: "simple line",
			args: args{
				loaders: map[LoaderName]Loader{
					"d": testLoader{},
					"b": testLoader{
						Order: Order{After: []LoaderName{"a"}, Before: []LoaderName{"c"}},
					},
					"a": testLoader{
						Order: Order{Before: []LoaderName{"b"}},
					},
					"c": testLoader{
						Order: Order{After: []LoaderName{"b", "a"}},
					},
					"e": testLoader{
						Order: Order{Before: []LoaderName{"c"}},
					},
					"f": testLoader{
						Order: Order{Before: []LoaderName{"e"}},
					},
				},
			},
			want: []LoaderName{"a", "f", "b", "e", "c", "d"},
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
