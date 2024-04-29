package defaultx

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/rakunlabs/chu/loader"
)

type Loader struct {
	hooks   []loader.HookFunc
	tagName string
}

func New(opts ...Option) *Loader {
	opt := &option{
		TagName: "default",
	}
	opt.apply(opts...)

	return &Loader{
		tagName: opt.TagName,
		hooks:   opt.Hooks,
	}
}

func (l Loader) Load(ctx context.Context, to any, opts ...loader.Option) error {
	v := reflect.ValueOf(to)
	if v.Kind() != reflect.Ptr {
		return errors.New("default: value is not a pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	if err := l.walk(ctx, v); err != nil {
		return fmt.Errorf("default: %w", err)
	}

	return nil
}

func (l *Loader) walk(ctx context.Context, v reflect.Value) error {
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
			tag := loader.TagValueM(fieldType, l.tagName)
			if tag == "-" {
				continue
			}

			if err := l.walkField(ctx, field, &tag); err != nil {
				return err
			}
		}
	default:
		return l.walkField(ctx, v, nil)
	}

	return nil
}

func (l *Loader) walkField(ctx context.Context, field reflect.Value, tag *string) error {
	// check direct exist
	if tag != nil && len(l.hooks) > 0 {
		var valGet interface{}
		var err error

		for _, hook := range l.hooks {
			valGet, err = hook(reflect.TypeFor[string](), field.Type(), tag)
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
		if err := l.walk(ctx, field); err != nil {
			return err
		}
	case reflect.Ptr:
		switch field.Type().Elem().Kind() {
		// just initialize the pointer if common type
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:

			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}

			if err := l.walkField(ctx, field.Elem(), tag); err != nil {
				return err
			}
		}
	default:
		if tag == nil {
			return nil
		}

		return loader.AssignValue(*tag, field)
	}

	return nil
}
