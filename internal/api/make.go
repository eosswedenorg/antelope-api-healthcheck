package api

import (
	"fmt"
)

func Make(name string, args ApiArguments) (ApiInterface, error) {
	factories := map[string]Factory{
		"v1":     AntelopeV1Factory,
		"v2":     AntelopeV2Factory,
		"atomic": AtomicAssetFactory,
		"debug":  DebugApiFactory,
	}

	if factory, ok := factories[name]; ok {
		return factory(args), nil
	}

	return nil, fmt.Errorf("invalid API '%s'", name)
}
