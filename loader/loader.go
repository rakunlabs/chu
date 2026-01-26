package loader

import "context"

type Loader interface {
	Load(ctx context.Context, to any, opt *Option) error
	LoadName() LoaderName
	LoadOrder() Order
}

type Order struct {
	Before []LoaderName
	After  []LoaderName
}

type LoaderName string

const (
	NameDefault LoaderName = "default"
	NameConsul  LoaderName = "consul"
	NameVault   LoaderName = "vault"
	NameHTTP    LoaderName = "http"
	NameFile    LoaderName = "file"
	NameEnv     LoaderName = "env"
)

var Loaders = map[LoaderName]Loader{}

func Add(l Loader) {
	Loaders[l.LoadName()] = l
}

var (
	OrderDefault = Order{
		Before: []LoaderName{NameConsul, NameVault, NameHTTP, NameFile, NameEnv},
	}
	OrderConsul = Order{
		Before: []LoaderName{NameVault, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault},
	}
	OrderVault = Order{
		Before: []LoaderName{NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul},
	}
	OrderHTTP = Order{
		Before: []LoaderName{NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault},
	}
	OrderFile = Order{
		Before: []LoaderName{NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault, NameHTTP},
	}
	OrderEnv = Order{
		After: []LoaderName{NameDefault, NameConsul, NameVault, NameHTTP, NameFile},
	}
)
