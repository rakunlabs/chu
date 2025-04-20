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
	NameHttp    = "http"
	NameFile    = "file"
	NameEnv     = "env"
)

var (
	DefaultOrderDefault = &Order{
		Before: []string{NameConsul, NameVault, NameHttp, NameFile, NameEnv},
	}
	DefaultOrderConsul = &Order{
		Before: []string{NameVault, NameHttp, NameFile, NameEnv},
		After:  []string{NameDefault},
	}
	DefaultOrderVault = &Order{
		Before: []string{NameHttp, NameFile, NameEnv},
		After:  []string{NameDefault, NameConsul},
	}
	DefaultOrderHttp = &Order{
		Before: []string{NameFile, NameEnv},
		After:  []string{NameDefault, NameConsul, NameVault},
	}
	DefaultOrderFile = &Order{
		Before: []string{NameEnv},
		After:  []string{NameDefault, NameConsul, NameVault, NameHttp},
	}
	DefaultOrderEnv = &Order{
		After: []string{NameDefault, NameConsul, NameVault, NameHttp, NameFile},
	}
)
