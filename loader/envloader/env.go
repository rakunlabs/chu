package envloader

import (
	"context"
	"errors"
	"os"
	"reflect"

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

func New(opts ...Option) *Loader {
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
	}
}

// Load loads the configuration from the environment.
//   - to must be a pointer to a struct
//   - only struct fields load values
func (l Loader) Load(ctx context.Context, to any) error {
	return l.LoadChu(ctx, to)
}

func (l Loader) LoadChu(ctx context.Context, to any, opts ...loader.Option) error {
	opt := loader.NewOption(opts...)

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
			tag = prefix + sanitizeTag(tag)

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
		var valGet interface{}
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
		// initialize slice if nil
		if field.IsNil() {
			maxValue := l.envValues.MaxValue(tag) + 1
			field.Set(reflect.MakeSlice(field.Type(), maxValue, maxValue))
		}

		for j := 0; j < field.Len(); j++ {
			if err := l.walk(ctx, field.Index(j), tag+"_"+cast.ToString(j)+"_"); err != nil {
				return err
			}
		}
	case reflect.Map:
		// skip map
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
