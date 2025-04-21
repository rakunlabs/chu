package defaultloader

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameDefault] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderDefault,
	}
}
