package loaderenv

import (
	"context"
	"errors"
	"os"
	"reflect"
	"strings"

	"github.com/rakunlabs/chu/loader"
	"github.com/spf13/cast"
)

type Loader struct {
	tagEnv    string
	tag       string
	envValues envHolder
	hooks     []loader.HookFunc
	envFiles  []string
	prefix    string
}

func New(opts ...Option) func() loader.Loader {
	return func() loader.Loader {
		opt := &option{
			TagEnv:   "env",
			Tag:      "cfg",
			EnvFiles: []string{".env"},
		}
		opt.apply(opts...)

		if envFile := os.Getenv("CONFIG_ENV_FILE"); envFile != "" {
			opt.EnvFiles = append(opt.EnvFiles, envFile)
		}

		return &Loader{
			envValues: opt.EnvHolder,
			hooks:     opt.Hooks,
			tagEnv:    opt.TagEnv,
			tag:       opt.Tag,
			envFiles:  opt.EnvFiles,
			prefix:    opt.Prefix,
		}
	}
}

// Load loads the configuration from the environment.
//   - to must be a pointer to a struct
//   - only struct fields load values
func (l Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	if len(opt.Hooks) > 0 {
		l.hooks = opt.Hooks
	}

	if opt.Tag != "" {
		l.tag = opt.Tag
	}

	v := reflect.ValueOf(to)
	if v.Kind() != reflect.Ptr {
		return errors.New("value is not a pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	envFileValues, err := getEnvValuesFromFiles(l.envFiles, l.prefix)
	if err != nil {
		return err
	}

	envValues := getEnvValues(l.prefix)
	for k, v := range envFileValues {
		envValues[k] = v
	}

	for k, v := range l.envValues {
		envValues[k] = v
	}

	l.envValues = envValues

	if err := l.walk(ctx, v, ""); err != nil {
		return err
	}

	return nil
}

func (l *Loader) walk(ctx context.Context, v reflect.Value, prefix string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// skip unexported field
			if !field.CanSet() {
				continue
			}

			fieldType := v.Type().Field(i)
			tag := loader.TagValue(fieldType, l.tagEnv, l.tag)
			if tag == "-" {
				continue
			}

			// always use uppercase tag
			if prefix != "" && !strings.HasSuffix(prefix, "_") {
				tag = prefix + "_" + sanitizeTag(tag)
			} else if prefix != "" {
				tag = prefix + sanitizeTag(tag)
			} else {
				tag = sanitizeTag(tag)
			}

			if err := l.walkField(ctx, field, tag); err != nil {
				return err
			}
		}
	default:
		field := v
		tag := prefix
		if tag != "" {
			tag = tag[:len(tag)-1]
		}

		return l.walkField(ctx, field, tag)
	}

	return nil
}

func (l *Loader) walkField(ctx context.Context, field reflect.Value, tag string) error {
	if !l.envValues.IsExist(tag) {
		return nil
	}

	// check direct exist
	if value, ok := l.envValues[tag]; ok && len(l.hooks) > 0 {
		var valGet any
		var err error

		for _, hook := range l.hooks {
			valGet, err = hook(reflect.TypeFor[string](), field.Type(), value)
			if err != nil {
				return err
			}
		}

		reflectValGet := reflect.ValueOf(valGet)
		if reflectValGet.Type() == field.Type() {
			field.Set(reflectValGet)
			return nil
		}
	}

	switch field.Kind() {
	case reflect.Struct:
		if err := l.walk(ctx, field, tag+"_"); err != nil {
			return err
		}
	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		if err := l.walk(ctx, field.Elem(), tag+"_"); err != nil {
			return err
		}
	case reflect.Slice:
		// Support slice of strings from env variable (e.g., "a,b,c")
		elKind := field.Type().Elem().Kind()
		value, ok := l.envValues[tag]
		if ok {
			strVal := cast.ToString(value)
			if strVal != "" {
				parts := splitAndTrim(strVal, ",")
				switch elKind {
				case reflect.String:
					field.Set(reflect.ValueOf(parts))
					return nil
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					intSlice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
					for i, p := range parts {
						v := cast.ToInt64(p)
						intSlice.Index(i).Set(reflect.ValueOf(v).Convert(field.Type().Elem()))
					}
					field.Set(intSlice)
					return nil
				case reflect.Float32, reflect.Float64:
					floatSlice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
					for i, p := range parts {
						v := cast.ToFloat64(p)
						floatSlice.Index(i).Set(reflect.ValueOf(v).Convert(field.Type().Elem()))
					}
					field.Set(floatSlice)
					return nil
				}
			}
		}
		// initialize slice if nil
		if field.IsNil() {
			maxValue := l.envValues.MaxValue(tag) + 1
			field.Set(reflect.MakeSlice(field.Type(), maxValue, maxValue))
		}
		for j := range field.Len() {
			if err := l.walk(ctx, field.Index(j), tag+"_"+cast.ToString(j)+"_"); err != nil {
				return err
			}
		}
	case reflect.Map:
		// Support map types from env variables (e.g., TEST_1_VALUE for tag TEST, key 1)
		keyType := field.Type().Key()
		prefix := tag + "_"
		result := reflect.MakeMap(field.Type())
		for k := range l.envValues {
			if strings.HasPrefix(k, prefix) {
				suffix := k[len(prefix):]
				mapKeyStr := suffix
				mapElemTag := k
				if idx := strings.Index(suffix, "_"); idx != -1 {
					mapKeyStr = suffix[:idx]
					mapElemTag = prefix + mapKeyStr + suffix[idx:]
				}
				var mapKey reflect.Value
				switch keyType.Kind() {
				case reflect.String:
					mapKey = reflect.ValueOf(mapKeyStr)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					mapKey = reflect.ValueOf(cast.ToInt64(mapKeyStr)).Convert(keyType)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					mapKey = reflect.ValueOf(cast.ToUint64(mapKeyStr)).Convert(keyType)
				case reflect.Bool:
					mapKey = reflect.ValueOf(cast.ToBool(mapKeyStr))
				case reflect.Float32, reflect.Float64:
					mapKey = reflect.ValueOf(cast.ToFloat64(mapKeyStr)).Convert(keyType)
				default:
					continue // skip unsupported key types
				}
				mapElem := reflect.New(field.Type().Elem()).Elem()
				if mapElem.Kind() == reflect.Struct {
					// For struct values, use prefix up to and including the underscore after the key
					structPrefix := mapElemTag
					if idx := strings.LastIndex(mapElemTag, "_"); idx != -1 {
						structPrefix = mapElemTag[:idx+1]
					}
					if err := l.walk(ctx, mapElem, structPrefix); err != nil {
						return err
					}
				} else {
					if err := l.walkField(ctx, mapElem, mapElemTag); err != nil {
						return err
					}
				}
				result.SetMapIndex(mapKey, mapElem)
			}
		}
		if result.Len() > 0 {
			field.Set(result)
		}
		return nil
	default:
		value, ok := l.envValues[tag]
		if !ok {
			return nil
		}

		return loader.AssignValue(value, field)
	}

	return nil
}

// Helper to split and trim string slices
func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
