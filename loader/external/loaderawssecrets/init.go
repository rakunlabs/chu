package loaderawssecrets

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Add(New())
}
