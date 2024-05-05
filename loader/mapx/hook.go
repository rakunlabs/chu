package mapx

import (
	"github.com/rakunlabs/chu/loader"

	"github.com/worldline-go/struct2"
)

func convertHookFuncs(hooks []loader.HookFunc) []struct2.HookDecodeFunc {
	if len(hooks) == 0 {
		return nil
	}

	hookFuncs := make([]struct2.HookDecodeFunc, len(hooks))
	for i, h := range hooks {
		hookFuncs[i] = struct2.HookDecodeFunc(h)
	}

	return hookFuncs
}
