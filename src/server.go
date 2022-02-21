package main

import (
    "fmt"
    "strings"
    "strconv"
    "internal/api"
    "github.com/eosswedenorg-go/eosapi"
    "github.com/eosswedenorg-go/haproxy"
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

    // var url string
    // var host string
    // var block_time int = 10
    // var version string = "v1"

    // Parse arguments.
    // -------------------
    split := strings.Split(strings.TrimSpace(args), "|")

    if len(split) < 2 {
        msg := "Invalid number of parameters in agent request"

        logger.Warn("Agent request error", "message", msg, "args", split)
        c.WriteString(fmt.Sprintf("%s#%s\n", haproxy.HealthCheckFailed, msg))
        c.Close()
        return
    }

    // 1. Api
    a.api = split[0]

    // 2. url (scheme + ip/domain + port)
    a.url = split[1]

    // 3. Block time.
    if len(split) > 1 {
        num, err := strconv.ParseInt(split[2], 10, 32)
        if err == nil {
            a.block_time = int(num)
        }
    }

    // 4. Host
    if len(split) > 3 {
        a.host = split[3]
    }

    // Check api.
    // -------------------
    healthCheckApi, err := createApi(&a)
    if err != nil {
        logger.Warn("Agent request error", "message", err)
        c.WriteString(fmt.Sprintf("%s#%s\n", haproxy.HealthCheckFailed, err))
        c.Close()
        return
    }

    status, msg := healthCheckApi.Call()

    logger.Info("API Check", append([]interface{}{
        "status", status},
        healthCheckApi.LogInfo().ToSlice()...)...)

    if status != haproxy.HealthCheckUp && len(msg) > 0 {
        logger.Warn("API Check Failed", "message", msg)
    }

    // Report status to HAproxy
    c.WriteString(fmt.Sprintln(status))
    c.Close()
}

//  spawnTcpServer
// ---------------------------------------------------------

func spawnTcpServer(addr string) {
    server := tcp_server.New(addr)
    server.OnMessage(onTcpMessage)
    server.Listen()
}
