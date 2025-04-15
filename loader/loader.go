package loader

import "context"

type Loader interface {
	LoadChu(ctx context.Context, to any, opt *Option) error
}
