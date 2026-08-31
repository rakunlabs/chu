package loaderdefault

import (
	"context"
	"errors"
	"reflect"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/gofret"
)

type Loader struct {
	hooks        []gofret.Hook
	codecOptions []gofret.Option
	codec        *gofret.Codec
	tagName      string
}

func New(opts ...Option) loader.Loader {
	opt := &option{
		TagName: "default",
	}
	opt.apply(opts...)

	return &Loader{
		tagName:      opt.TagName,
		hooks:        opt.Hooks,
		codecOptions: opt.CodecOptions,
	}
}

func (l *Loader) LoadName() loader.LoaderName {
	return loader.NameDefault
}

func (l *Loader) LoadOrder() loader.Order {
	return loader.OrderDefault
}

func (l Loader) Load(ctx context.Context, to any, opt *loader.Option) error {
	if len(opt.Hooks) > 0 {
		l.hooks = opt.Hooks
	}
	if len(opt.CodecOptions) > 0 {
		l.codecOptions = opt.CodecOptions
	}

	codecOptions := append([]gofret.Option{}, l.codecOptions...)
	codecOptions = append(codecOptions, gofret.WithHooks(l.hooks...))
	l.codec = gofret.New(codecOptions...)

	v := reflect.ValueOf(to)
	if v.Kind() != reflect.Ptr {
		return errors.New("value is not a pointer")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil
	}

	if err := l.walk(ctx, v); err != nil {
		return err
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
		for i := range v.NumField() {
			field := v.Field(i)
			// skip unexported field
			if !field.CanSet() {
				continue
			}

			fieldType := v.Type().Field(i)
			tag := loader.TagValueM(fieldType, l.tagName)
			if tag != nil && *tag == "-" {
				continue
			}

			if err := l.walkField(ctx, field, tag); err != nil {
				return err
			}
		}
	default:
		return l.walkField(ctx, v, nil)
	}

	return nil
}

func (l *Loader) walkField(ctx context.Context, field reflect.Value, tag *string) error {
	switch field.Kind() {
	case reflect.Struct:
		if err := l.walk(ctx, field); err != nil {
			return err
		}
	case reflect.Ptr:
		if tag == nil {
			return nil
		}

		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		if err := l.walkField(ctx, field.Elem(), tag); err != nil {
			return err
		}
	default:
		if tag == nil {
			return nil
		}

		return l.codec.ToInto(*tag, field.Addr().Interface())
	}

	return nil
}
