
package main

import (
    "strings"
    "fmt"
    "strconv"
    "github.com/eosswedenorg/eosio-api-healthcheck/src/api"
)

type arguments struct {
    url string
    host string
    num_blocks int
}

func ParseArguments(args []string) arguments {

    a := arguments{
        num_blocks: 10,
    }

    // 1. url (scheme + ip/domain + port)
    a.url = args[0]

    // 2. num blocks
    if len(args) > 1 {
        num, err := strconv.ParseInt(args[1], 10, 32)
        if err == nil {
            a.num_blocks = int(num)
        }
    }

    // 3. Host
    if len(args) > 2 {
        a.host = args[2]
    }

    return a
}

func ParseRequest(request string) (api.ApiInterface, error) {

    // Parse arguments.
    // -------------------
    p := strings.Split(strings.TrimSpace(request), "|")

    if len(p) < 2 {
        return nil, fmt.Errorf("Invalid number of parameters in agent request")
    }

    a := ParseArguments(p[1:])

    switch p[0] {
    case "v1":
        return api.NewEosioV1(a.url, a.host, float64(a.num_blocks / 2)), nil
    case "v2":
        return api.NewEosioV2(a.url, a.host, int64(a.num_blocks)), nil
    case "contract":
        return api.NewEosioContract(a.url, float64(a.num_blocks / 2)), nil
    case "debug":
        return api.NewDebugApi(a.url), nil
    }

    return nil, fmt.Errorf("Invalid API '%s'", p[0])
}
