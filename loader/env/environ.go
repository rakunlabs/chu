package env

import (
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-envparse"
	"github.com/spf13/cast"
)

type envHolder map[string]string

var rgxIndex = regexp.MustCompile(`^_(\d+)_?.*$`)

func regexIndex(key string) int {
	matches := rgxIndex.FindStringSubmatch(key)
	if len(matches) != 2 {
		return -1
	}

	return cast.ToInt(matches[1])
}

func (e envHolder) IsExist(key string) bool {
	for k := range e {
		if strings.HasPrefix(k, key) {
			return true
		}
	}

	return false
}

func (e envHolder) MaxValue(key string) int {
	max := 0
	for k := range e {
		if strings.HasPrefix(k, key) {
			idx := regexIndex(strings.TrimPrefix(k, key))
			if idx > max {
				max = idx
			}
		}
	}

	return max
}

func getEnvValues(prefix string) envHolder {
	env := os.Environ()
	values := make(map[string]string, len(env))
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		pair[0] = sanitizePrefix(prefix, pair[0])
		if pair[0] != "" {
			values[pair[0]] = pair[1]
		}
	}

	return values
}

func getEnvValuesFromFiles(files []string, prefix string) (envHolder, error) {
	values := make(map[string]string)
	for _, file := range files {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}

		file, err := os.Open(file)
		if err != nil {
			return nil, err
		}

		env, err := envparse.Parse(file)
		for k, v := range env {
			k = sanitizePrefix(prefix, k)
			if k != "" {
				values[k] = v
			}
		}
	}

	return values, nil
}

func sanitizeTag(tag string) string {
	return strings.ReplaceAll(strings.ToUpper(tag), "-", "_")
}

func sanitizePrefix(prefix string, key string) string {
	key = sanitizeTag(key)
	prefix = sanitizeTag(prefix)

	if prefix == "" {
		return key
	}

	return strings.TrimPrefix(key, prefix)
}
