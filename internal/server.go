package internal

import (
    "strings"
    log "github.com/inconshreveable/log15"
    "github.com/eosswedenorg/eosio-api-healthcheck/internal/api"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
    "github.com/eosswedenorg-go/tcp_server"
)

//  onTcpMessage callback function
// ---------------------------------------------------------

func onTcpMessage(c *tcp_server.Client, args string) {

    logger := log.Root()

    // Check api.
    // -------------------
    healthCheckApi, err := ParseRequest(args)
    if err != nil {
        logger.Warn("Agent request error", "message", err)
        resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")

        c.WriteString(resp.String())
        c.Close()
        return
    }

    status, msg := healthCheckApi.Call()

    params := api.LogParams{}
    params.Add("status", strings.TrimSpace(status.String()))

    if msg != "OK" && len(msg) > 0 {
        params.Add("error", msg)
    }

    logger.Info("API Check", params.Combine(healthCheckApi.LogInfo())...)
    // Report status to HAproxy
    c.WriteString(status.String())
    c.Close()
}

//  SpawnTcpServer
// ---------------------------------------------------------

func SpawnTcpServer(addr string) error {
    server := tcp_server.New(addr)
    server.OnMessage(onTcpMessage)

    err := server.Connect()
    if err == nil {
        go server.Listen()
    }
    return err
}
