package loader

import (
	"reflect"

	"github.com/spf13/cast"
)

// TagValue returns the value of the tag in the field in order of the tags.
//
// If the tag is not found, it will return the field name.
func TagValue(field reflect.StructField, tags ...string) string {
	value := field.Name
	if len(tags) == 0 {
		return value
	}

	for _, tag := range tags {
		if tag == "" {
			continue
		}

		if v := field.Tag.Get(tag); v != "" {
			return v
		}
	}

	return value
}

// TagValueM returns the value of the tag in the field in order of the tags.
//
// If the tag is not found, it will return an empty string.
func TagValueM(field reflect.StructField, tags ...string) string {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		if v, ok := field.Tag.Lookup(tag); ok {
			return v
		}
	}

	return ""
}

// AssignValue assigns the value to the field.
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
		return nil
	}

	return nil
}
