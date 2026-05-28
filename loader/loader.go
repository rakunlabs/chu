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
	NameDefault      LoaderName = "default"
	NameConsul       LoaderName = "consul"
	NameVault        LoaderName = "vault"
	NameGCPParameter LoaderName = "gcpparameter"
	NameGCPSecret    LoaderName = "gcpsecret"
	NameHTTP         LoaderName = "http"
	NameFile         LoaderName = "file"
	NameEnv          LoaderName = "env"
)

var Loaders = map[LoaderName]Loader{}

func Add(l Loader) {
	Loaders[l.LoadName()] = l
}

var (
	OrderDefault = Order{
		Before: []LoaderName{NameConsul, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
	}
	OrderConsul = Order{
		Before: []LoaderName{NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault},
	}
	OrderVault = Order{
		Before: []LoaderName{NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul},
	}
	OrderGCPParameter = Order{
		Before: []LoaderName{NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault},
	}
	OrderGCPSecret = Order{
		Before: []LoaderName{NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault, NameGCPParameter},
	}
	OrderHTTP = Order{
		Before: []LoaderName{NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault, NameGCPParameter, NameGCPSecret},
	}
	OrderFile = Order{
		Before: []LoaderName{NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP},
	}
	OrderEnv = Order{
		After: []LoaderName{NameDefault, NameConsul, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile},
	}
)
