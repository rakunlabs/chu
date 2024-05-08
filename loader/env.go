package loader

import "os"

func GetExistEnv(name ...string) (string, bool) {
	for _, n := range name {
		if v, ok := os.LookupEnv(n); ok {
			return v, true
		}
	}

	return "", false
}
