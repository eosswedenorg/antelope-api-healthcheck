package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eosswedenorg/antelope-api-healthcheck/internal/api"
)

func ParseArguments(args []string) api.ApiArguments {
	a := api.ApiArguments{
		NumBlocks: 10,
	}

	// 1. url (scheme + ip/domain + port)
	a.Url = args[0]

	// 2. num blocks
	if len(args) > 1 {
		num, err := strconv.ParseInt(args[1], 10, 32)
		if err == nil {
			a.NumBlocks = int(num)
		}
	}

	// 3. Host
	if len(args) > 2 {
		a.Host = args[2]
	}

	return a
}

func ParseRequest(request string) (api.ApiInterface, error) {
	p := strings.Split(strings.TrimSpace(request), "|")

	if len(p) < 2 {
		return nil, fmt.Errorf("invalid number of parameters in agent request")
	}

	a := ParseArguments(p[1:])

	return api.Make(p[0], a)
}
