package main

import (
    "strings"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
    "github.com/eosswedenorg-go/tcp_server"
)

//  onTcpMessage callback function
// ---------------------------------------------------------

func onTcpMessage(c *tcp_server.Client, args string) {

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
