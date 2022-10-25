
package main

import (
    "strings"
    "fmt"
    "strconv"
    "github.com/eosswedenorg/eosio-api-healthcheck/src/api"
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

    // Parse arguments.
    // -------------------
    p := strings.Split(strings.TrimSpace(request), "|")

    if len(p) < 2 {
        return nil, fmt.Errorf("invalid number of parameters in agent request")
    }

    a := ParseArguments(p[1:])

    switch p[0] {
    case "v1":
        return api.NewEosioV1(a.Url, a.Host, float64(a.NumBlocks / 2)), nil
    case "v2":
        return api.NewEosioV2(a.Url, a.Host, int64(a.NumBlocks)), nil
    case "contract":
        return api.NewEosioContract(a.Url, float64(a.NumBlocks / 2)), nil
    case "debug":
        return api.NewDebugApi(a.Url), nil
    }

    return nil, fmt.Errorf("invalid API '%s'", p[0])
}
