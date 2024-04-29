package loader

import (
	"fmt"
	"reflect"

	"github.com/spf13/cast"
)

var TagName = "cfg"

func TagValue(field reflect.StructField, tags ...string) string {
	value := field.Name
	if len(tags) == 0 {
		return value
	}

	for _, tag := range tags {
		if v := field.Tag.Get(tag); v != "" {
			return v
		}
	}

	return value
}

func TagValueM(field reflect.StructField, tags ...string) string {
	for _, tag := range tags {
		if v, ok := field.Tag.Lookup(tag); ok {
			return v
		}
	}

	return ""
}

func AssignValue(value string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.SetInt(cast.ToInt64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.SetUint(cast.ToUint64(value))
	case reflect.Float32, reflect.Float64:
		field.SetFloat(cast.ToFloat64(value))
	case reflect.Bool:
		field.SetBool(cast.ToBool(value))
	default:
		return fmt.Errorf("unsupported type %s", field.Kind())
	}

	return nil
}
