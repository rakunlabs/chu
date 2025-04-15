package fileloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rakunlabs/chu/loader"
)

type Loader struct {
	FileSuffix []string
	Folders    []string
	Decoders   map[string]Decoder
	MapDecoder func(data any, to any) error
}

var LoaderName = "file"

func New(opts ...Option) func() loader.Loader {
	return func() loader.Loader {
		opt := &option{
			FileSuffix: []string{".toml", ".yaml", ".yml", ".json"},
			Folders:    []string{"/etc"},
			Decoders:   Decoders(),
		}
		opt.apply(opts...)

		return &Loader{
			FileSuffix: opt.FileSuffix,
			Folders:    opt.Folders,
			Decoders:   opt.Decoders,
		}
	}
}

// Load loads the configuration from the file.
//   - first it checks the current directory after that it checks the etc folder.
//   - CONFIG_PATH environment variable is used to determine the file path.
func (l Loader) LoadChu(ctx context.Context, to any, opt *loader.Option) error {
	if l.MapDecoder == nil {
		l.MapDecoder = opt.MapDecoder
	}

	if l.MapDecoder == nil {
		return fmt.Errorf("map decoder is not set %w", loader.ErrMissingOpt)
	}

	if path := l.getEnv(opt.Name); path != "" {
		if err := l.loadTo(ctx, path, to); err != nil {
			return err
		}

		return nil
	}

	if opt.Name == "" {
		return nil
	}

	path := l.getPath(opt.Name)
	if path == "" {
		return nil
	}

	if err := l.loadTo(ctx, path, to); err != nil {
		return err
	}

	return nil
}

func (l Loader) getEnv(name string) string {
	if path := os.Getenv("CONFIG_PATH" + "_" + strings.ToUpper(name)); path != "" {
		return path
	}

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}

	return ""
}

func (l Loader) getPath(name string) string {
	// check current directory
	for _, suffix := range l.FileSuffix {
		path := name + suffix
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// check other folder
	for _, folder := range l.Folders {
		for _, suffix := range l.FileSuffix {
			path := filepath.Join(folder, name+suffix)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

func (l Loader) loadTo(_ context.Context, path string, to any) error {
	mapping, err := l.fileToMap(path)
	if err != nil {
		return err
	}

	return l.MapDecoder(mapping, to)
}

func (l Loader) fileToMap(path string) (any, error) {
	fileDecoder, err := l.getFileDecoder(filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	var mapping any
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := fileDecoder.Decode(file, &mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
