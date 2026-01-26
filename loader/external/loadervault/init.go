package loadervault

import "github.com/rakunlabs/chu/loader"

func init() {
	loader.Add(New())
}
