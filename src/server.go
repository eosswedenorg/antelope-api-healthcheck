package main

import (
    "fmt"
    "strings"
    "strconv"
    "github.com/eosswedenorg/eosio-api-healthcheck/src/utils"
    "github.com/eosswedenorg/eosio-api-healthcheck/src/api"
    "github.com/eosswedenorg-go/eosapi"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
    "github.com/eosswedenorg-go/tcp_server"
)

type arguments struct {
    api string
    url string
    host string
    block_time int
}

func createApi(a *arguments) (api.ApiInterface, error) {

    switch a.api {
    case "v1":
        return api.NewEosioV1(eosapi.ReqParams{Url: a.url, Host: a.host}, float64(a.block_time)), nil
    case "v2":
        return api.NewEosioV2(eosapi.ReqParams{Url: a.url, Host: a.host}, int64(a.block_time / 2)), nil
    case "contract":
        return api.NewEosioContract(a.url, float64(a.block_time)), nil
    }

    return nil, fmt.Errorf("Invalid API '%s'", a.api)
}

//  onTcpMessage callback function
// ---------------------------------------------------------

func onTcpMessage(c *tcp_server.Client, args string) {

    a := arguments{
        api: "v1",
        block_time: 10,
    }

    // Parse arguments.
    // -------------------
    split := strings.Split(strings.TrimSpace(args), "|")

    if len(split) < 2 {
        msg := "Invalid number of parameters in agent request"

        logger.Warn("Agent request error", "message", msg, "args", split)
        resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, msg)

        c.WriteString(resp.String())
        c.Close()
        return
    }

    // Old format: <url> <version> <api> <block_time>
    if utils.IsUrl(split[0]) {

        logger.Warn("Deprecated format. Please change to the new format: <api>|<url>[|<block_time>|<host>]")

        // 1. url (scheme + ip/domain + port)
        a.url = split[0]

        // 2. api
        if len(split) > 1 {
            a.api = split[1]
        }

        // 3. Block time
        if len(split) > 2 {
            num, err := strconv.ParseInt(split[2], 10, 32)
            if err == nil {
                a.block_time = int(num)
            }
        }

        // 4. Host
        if len(split) > 3 {
            a.host = split[3]
        }

    } else {

        if len(split) < 2 {
            msg := "Invalid number of parameters in agent request"

            logger.Warn("Agent request error", "message", msg, "args", split)
            resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, msg)

            c.WriteString(resp.String())
            c.Close()
            return
        }

        // 1. Api
        a.api = split[0]

        // 2. url (scheme + ip/domain + port)
        a.url = split[1]

        // 3. Block time.
        if len(split) > 2 {
            num, err := strconv.ParseInt(split[2], 10, 32)
            if err == nil {
                a.block_time = int(num)
            }
        }

        // 4. Host
        if len(split) > 3 {
            a.host = split[3]
        }
    }

    // Check api.
    // -------------------
    healthCheckApi, err := createApi(&a)
    if err != nil {
        logger.Warn("Agent request error", "message", err)
        resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, err.Error())

        c.WriteString(resp.String())
        c.Close()
        return
    }

    status, msg := healthCheckApi.Call()

    logger.Info("API Check", append([]interface{}{
        "status", strings.TrimSpace(status.String())},
        healthCheckApi.LogInfo().ToSlice()...)...)

    if msg != "OK" && len(msg) > 0 {
        logger.Warn("API Check Failed", "message", msg)
    }

    // Report status to HAproxy
    c.WriteString(status.String())
    c.Close()
}

//  spawnTcpServer
// ---------------------------------------------------------

func spawnTcpServer(addr string) {
    server := tcp_server.New(addr)
    server.OnMessage(onTcpMessage)
    server.Listen()
}
