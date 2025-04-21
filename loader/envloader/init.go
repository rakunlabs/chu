package envloader

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameEnv] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderEnv,
	}
}
