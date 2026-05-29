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

func Test_orderLoaders_canonicalChain(t *testing.T) {
	loaders := map[LoaderName]Loader{
		NameDefault:       testLoader{Name: NameDefault, Order: OrderDefault},
		NameConsul:        testLoader{Name: NameConsul, Order: OrderConsul},
		NameAWSSSM:        testLoader{Name: NameAWSSSM, Order: OrderAWSSSM},
		NameAWSSecrets:    testLoader{Name: NameAWSSecrets, Order: OrderAWSSecrets},
		NameAzureKeyVault: testLoader{Name: NameAzureKeyVault, Order: OrderAzureKeyVault},
		NameVault:         testLoader{Name: NameVault, Order: OrderVault},
		NameGCPParameter:  testLoader{Name: NameGCPParameter, Order: OrderGCPParameter},
		NameGCPSecret:     testLoader{Name: NameGCPSecret, Order: OrderGCPSecret},
		NameHTTP:          testLoader{Name: NameHTTP, Order: OrderHTTP},
		NameFile:          testLoader{Name: NameFile, Order: OrderFile},
		NameEnv:           testLoader{Name: NameEnv, Order: OrderEnv},
	}

	want := []LoaderName{
		NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault,
		NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv,
	}

	if got := OrderLoaders(loaders); !reflect.DeepEqual(got, want) {
		t.Errorf("OrderLoaders() = %v, want %v", got, want)
	}
}
