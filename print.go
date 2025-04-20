package chu

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/rakunlabs/chu/loader"
	"github.com/spf13/cast"
)

// Print is a function that takes a context and an interface{} value,
// and returns a JSON representation of the value.
//   - Uses "log" tag and "-" to skip fields or false to skip
func PrintE(ctx context.Context, v any) (string, error) {
	m, err := buildLoggableMap(ctx, reflect.ValueOf(v))
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func Print(ctx context.Context, v any) string {
	result, _ := PrintE(ctx, v)

	return result
}

// buildLoggableMap recursively builds a map representation of v, skipping fields with log:"false" or log:"-".
func buildLoggableMap(ctx context.Context, v reflect.Value) (any, error) {
	if !v.IsValid() {
		return nil, nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil
		}
		return buildLoggableMap(ctx, v.Elem())
	}
	if v.Kind() == reflect.Struct {
		if newV, ok := overrideValue(v); ok {
			return newV, nil
		}

		m := make(map[string]any)
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			fieldType := t.Field(i)
			// skip unexported fields
			if fieldType.PkgPath != "" {
				continue
			}
			// skip unsupported kinds
			if field.Kind() == reflect.Func || field.Kind() == reflect.Chan || field.Kind() == reflect.UnsafePointer || field.Kind() == reflect.Uintptr || field.Kind() == reflect.Complex64 || field.Kind() == reflect.Complex128 {
				continue
			}
			tag := loader.TagValueM(fieldType, "log")
			if tag != nil {
				if *tag == "-" {
					continue
				}
				if v, _ := strconv.ParseBool(*tag); !v {
					continue
				}
			}
			key := loader.TagValue(fieldType, "cfg")
			val, err := buildLoggableMap(ctx, field)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
		return m, nil
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		arr := make([]any, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			if item.Kind() == reflect.Func || item.Kind() == reflect.Chan || item.Kind() == reflect.UnsafePointer {
				continue
			}
			val, err := buildLoggableMap(ctx, item)
			if err != nil {
				return nil, err
			}
			arr = append(arr, val)
		}
		return arr, nil
	}
	if v.Kind() == reflect.Map {
		m := make(map[string]any)
		for _, key := range v.MapKeys() {
			keyStr := cast.ToString(key.Interface())
			if keyStr == "" {
				continue
			}

			val := v.MapIndex(key)
			if val.Kind() == reflect.Func || val.Kind() == reflect.Chan || val.Kind() == reflect.UnsafePointer || val.Kind() == reflect.Uintptr || val.Kind() == reflect.Complex64 || val.Kind() == reflect.Complex128 {
				continue
			}
			mappedVal, err := buildLoggableMap(ctx, val)
			if err != nil {
				return nil, err
			}
			m[keyStr] = mappedVal
		}
		return m, nil
	}

	// For other types, return the value as interface{}
	if newV, ok := overrideValue(v); ok {
		return newV, nil
	}

	return v.Interface(), nil
}

func overrideValue(v reflect.Value) (any, bool) {
	// For other types, return the value as interface{}
	if v.Type() == durationReflectType {
		return printDuration(v.Interface().(time.Duration)), true
	}

	return nil, false
}

var durationReflectType = reflect.TypeOf(time.Duration(0))

type printDuration time.Duration

func (d printDuration) MarshalJSON() ([]byte, error) {
	// Convert the duration to a string representation
	durationStr := time.Duration(d).String()
	// Marshal the string representation to JSON
	return json.Marshal(durationStr)
}
