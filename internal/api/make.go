package api

import (
	"fmt"
)

func Make(name string, args ApiArguments) (ApiInterface, error) {
	factories := map[string]Factory{
		"v1":     EosioV1Factory,
		"v2":     EosioV2Factory,
		"atomic": AtomicAssetFactory,
		"debug":  DebugApiFactory,
	}

	if factory, ok := factories[name]; ok {
		return factory(args), nil
	}

	return nil, fmt.Errorf("invalid API '%s'", name)
}
