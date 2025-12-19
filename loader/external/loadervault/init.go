package loadervault

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameVault] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderVault,
	}
}
