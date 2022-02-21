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

func createApi(a *arguments) api.ApiInterface {

    switch a.api {
    case "v2":
        return api.NewEosioV2(eosapi.ReqParams{Url: a.url, Host: a.host}, int64(a.block_time / 2))
    }

    return api.NewEosioV1(eosapi.ReqParams{Url: a.url, Host: a.host}, float64(a.block_time))
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

    // 1. url (scheme + ip/domain + port)
    a.url = split[0]

    // 2. Block time.
    if len(split) > 1 {
        num, err := strconv.ParseInt(split[1], 10, 32)
        if err == nil {
            a.block_time = int(num)
        }
    }

    // 3. Api
    if len(split) > 2 {
        a.api = split[2]
    }

    // 4. Host
    if len(split) > 3 {
        a.host = split[3]
    }

    // Check api.
    // -------------------
    healthCheckApi := createApi(&a)

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
