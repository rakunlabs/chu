package loader

import (
	"reflect"
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
func TagValueM(field reflect.StructField, tags ...string) *string {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		if v, ok := field.Tag.Lookup(tag); ok {
			return &v
		}
	}

	return nil
}
