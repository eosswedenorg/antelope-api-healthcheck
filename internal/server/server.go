package server

import (
	"strings"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg-go/tcp_server"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/api"
	log "github.com/inconshreveable/log15"
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
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")

		_, err = c.WriteString(resp.String())
		if err != nil {
			logger.Error("WriteString", "message", err)
		}

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
	_, err = c.WriteString(status.String())
	if err != nil {
		logger.Error("WriteString", "message", err)
	}
	c.Close()
}

//  Start
// ---------------------------------------------------------

func Start(addr string) (*tcp_server.Server, error) {
	server := tcp_server.New(addr)
	server.OnMessage(onTcpMessage)

	err := server.Connect()
	if err == nil {
		err = server.Listen()
	}
	return server, err
}
