package loader

import (
	"reflect"
	"time"

	"github.com/spf13/cast"
	"github.com/xhit/go-str2duration/v2"
)

// HookFunc get input, output and data and return modified data.
type HookFunc func(reflect.Type, reflect.Type, interface{}) (interface{}, error)

// HookTimeDuration for time.Duration
func HookTimeDuration(in reflect.Type, out reflect.Type, data interface{}) (interface{}, error) {
	if out == reflect.TypeFor[time.Duration]() {
		switch in.Kind() {
		case reflect.String:
			return str2duration.ParseDuration(data.(string))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return time.Duration(cast.ToInt64(data)), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return time.Duration(cast.ToUint64(data)), nil
		case reflect.Float32, reflect.Float64:
			return time.Duration(cast.ToFloat64(data)), nil
		}
	}

	return data, nil
}
