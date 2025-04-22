package loaderfile

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Loaders[loader.NameFile] = loader.LoadHolder{
		Loader: New(),
		Order:  loader.OrderFile,
	}
}
