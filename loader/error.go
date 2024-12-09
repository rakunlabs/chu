package loader

import "errors"

var (
	ErrSkipLoader = errors.New("skip loader")
	ErrMissingOpt = errors.New("missing option")
)
