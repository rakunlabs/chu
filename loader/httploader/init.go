package httploader

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameHTTP] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderHTTP,
	}
}
