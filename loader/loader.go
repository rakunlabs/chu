package loader

import "context"

type Loader interface {
	LoadChu(ctx context.Context, to any, opt *Option) error
}

type LoadHolder struct {
	Loader func() Loader
	Order  *Order
}

type Order struct {
	Before []string
	After  []string
}

const (
	NameDefault = "default"
	NameConsul  = "consul"
	NameVault   = "vault"
	NameHTTP    = "http"
	NameFile    = "file"
	NameEnv     = "env"
)

var Loaders = map[string]LoadHolder{}

var (
	OrderDefault = &Order{
		Before: []string{NameConsul, NameVault, NameHTTP, NameFile, NameEnv},
	}
	OrderConsul = &Order{
		Before: []string{NameVault, NameHTTP, NameFile, NameEnv},
		After:  []string{NameDefault},
	}
	OrderVault = &Order{
		Before: []string{NameHTTP, NameFile, NameEnv},
		After:  []string{NameDefault, NameConsul},
	}
	OrderHTTP = &Order{
		Before: []string{NameFile, NameEnv},
		After:  []string{NameDefault, NameConsul, NameVault},
	}
	OrderFile = &Order{
		Before: []string{NameEnv},
		After:  []string{NameDefault, NameConsul, NameVault, NameHTTP},
	}
	OrderEnv = &Order{
		After: []string{NameDefault, NameConsul, NameVault, NameHTTP, NameFile},
	}
)
