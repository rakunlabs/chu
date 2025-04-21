package consulloader

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameConsul] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderConsul,
	}
}
