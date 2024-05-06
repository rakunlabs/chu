package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rakunlabs/chu/loader"
	"github.com/rakunlabs/chu/loader/mapx"
)

type Loader struct {
	mapx                  *mapx.Loader
	hooks                 []loader.HookFunc
	weaklyIgnoreSeperator bool
	weaklyDashUnderscore  bool
	fileSuffix            []string
	folders               []string
	name                  string
	decoders              map[string]Decoder
}

func New(opts ...Option) *Loader {
	opt := &option{
		FileSuffix:            []string{".toml", ".yaml", ".yml", ".json"},
		Decoders:              getDecoders(),
		Folders:               []string{"/etc"},
		WeaklyIgnoreSeperator: true,
		WeaklyDashUnderscore:  false,
	}
	opt.apply(opts...)

	return &Loader{
		hooks:                 opt.Hooks,
		weaklyIgnoreSeperator: opt.WeaklyIgnoreSeperator,
		weaklyDashUnderscore:  opt.WeaklyDashUnderscore,
		fileSuffix:            opt.FileSuffix,
		folders:               opt.Folders,
		decoders:              opt.Decoders,
	}
}

// Load loads the configuration from the file.
//   - first it checks the current directory after that it checks the etc folder.
//   - CONFIG_PATH environment variable is used to determine the file path.
func (l Loader) Load(ctx context.Context, name string, to any) error {
	return l.LoadChu(ctx, to, loader.WithName(name))
}

func (l Loader) LoadChu(ctx context.Context, to any, opts ...loader.Option) error {
	opt := loader.NewOption(opts...)

	if opt.Name != "" {
		l.name = opt.Name
	}

	if len(opt.Hooks) > 0 {
		l.hooks = opt.Hooks
	}

	l.mapx = mapx.New(
		mapx.WithTag(opt.Tag),
		mapx.WithHooks(l.hooks...),
		mapx.WithWeaklyIgnoreSeperator(l.weaklyIgnoreSeperator),
		mapx.WithWeaklyDashUnderscore(l.weaklyDashUnderscore),
	)

	if path := os.Getenv("CONFIG_PATH"); path != "" {
		if err := l.loadTo(ctx, path, to); err != nil {
			return fmt.Errorf("file: %w", err)
		}

		return nil
	}

	if l.name == "" {
		return nil
	}

	path := l.getPath(l.name)
	if path == "" {
		return nil
	}

	if err := l.loadTo(ctx, path, to); err != nil {
		return fmt.Errorf("file: %w", err)
	}

	return nil
}

func (l Loader) getPath(name string) string {
	// check current directory
	for _, suffix := range l.fileSuffix {
		path := filepath.Join(name, suffix)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// check other folder
	for _, folder := range l.folders {
		for _, suffix := range l.fileSuffix {
			path := filepath.Join(folder, name, suffix)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	return ""
}

func (l Loader) loadTo(ctx context.Context, path string, to any) error {
	mapping, err := l.fileToMap(path)
	if err != nil {
		return err
	}

	return l.mapx.SetValue(mapping).Load(ctx, to)
}

func (l Loader) fileToMap(path string) (interface{}, error) {
	decoder, err := l.getFileDecoder(filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	var mapping interface{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(file, &mapping); err != nil {
		return nil, err
	}

	return mapping, nil
}
