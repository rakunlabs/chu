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
	NameDefault       LoaderName = "default"
	NameConsul        LoaderName = "consul"
	NameAWSSSM        LoaderName = "awsssm"
	NameAWSSecrets    LoaderName = "awssecrets"
	NameAzureKeyVault LoaderName = "azurekeyvault"
	NameVault         LoaderName = "vault"
	NameGCPParameter  LoaderName = "gcpparameter"
	NameGCPSecret     LoaderName = "gcpsecret"
	NameHTTP          LoaderName = "http"
	NameFile          LoaderName = "file"
	NameEnv           LoaderName = "env"
)

var Loaders = map[LoaderName]Loader{}

func Add(l Loader) {
	Loaders[l.LoadName()] = l
}

// Canonical load order (lowest to highest precedence; later loaders override
// earlier ones):
//
//	default -> consul -> awsssm -> awssecrets -> azurekeyvault -> vault ->
//	gcpparameter -> gcpsecret -> http -> file -> env
//
// Each Order lists every known loader so the chain stays consistent even when
// some loaders are not imported (OrderLoaders skips names that are absent).
var (
	OrderDefault = Order{
		Before: []LoaderName{NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
	}
	OrderConsul = Order{
		Before: []LoaderName{NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault},
	}
	OrderAWSSSM = Order{
		Before: []LoaderName{NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul},
	}
	OrderAWSSecrets = Order{
		Before: []LoaderName{NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM},
	}
	OrderAzureKeyVault = Order{
		Before: []LoaderName{NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets},
	}
	OrderVault = Order{
		Before: []LoaderName{NameGCPParameter, NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault},
	}
	OrderGCPParameter = Order{
		Before: []LoaderName{NameGCPSecret, NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault},
	}
	OrderGCPSecret = Order{
		Before: []LoaderName{NameHTTP, NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter},
	}
	OrderHTTP = Order{
		Before: []LoaderName{NameFile, NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret},
	}
	OrderFile = Order{
		Before: []LoaderName{NameEnv},
		After:  []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP},
	}
	OrderEnv = Order{
		After: []LoaderName{NameDefault, NameConsul, NameAWSSSM, NameAWSSecrets, NameAzureKeyVault, NameVault, NameGCPParameter, NameGCPSecret, NameHTTP, NameFile},
	}
)
