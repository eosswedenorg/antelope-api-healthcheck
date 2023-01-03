package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/api"
	log "github.com/inconshreveable/log15"
	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	gnet.BuiltinEventEngine

	addr string
	eng  gnet.Engine
}

func New(addr string) *Server {
	return &Server{
		addr: fmt.Sprintf("tcp://%s", addr),
	}
}

//	OnBoot callback function
//
// ---------------------------------------------------------
func (s *Server) OnBoot(eng gnet.Engine) gnet.Action {
	s.eng = eng
	log.Info("Server started", "addr", s.addr)
	return gnet.None
}

//	OnShutdown callback function
//
// ---------------------------------------------------------
func (s *Server) OnShutdown(eng gnet.Engine) {
	log.Info("Server shutdown")
}

//	OnTick callback function
//
// ---------------------------------------------------------
func (s *Server) OnTick() (time.Duration, gnet.Action) {
	log.Info("Server info", "connections", s.eng.CountConnections())
	return time.Second * 10, gnet.None
}

//	OnTraffic callback function
//
// ---------------------------------------------------------
func (s *Server) OnTraffic(c gnet.Conn) gnet.Action {
	logger := log.Root()

	req, err := c.Next(-1)
	if err != nil {
		logger.Error("Read", "message", err)
		return gnet.Close
	}

	// Check api.
	// -------------------
	healthCheckApi, err := ParseRequest(string(req))
	if err != nil {
		logger.Warn("Agent request error", "message", err)
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")

		_, err = c.Write([]byte(resp.String()))
		if err != nil {
			logger.Error("Write", "message", err)
		}

		return gnet.Close
	}

	// gnet library does not like blocking calls.
	// as we do a blocking http call here, we need to wrap it in a goroutine.
	go func() {
		status, msg := healthCheckApi.Call()

		params := api.LogParams{}
		params.Add("status", strings.TrimSpace(status.String()))

		if msg != "OK" && len(msg) > 0 {
			params.Add("error", msg)
		}

		logger.Info("API Check", params.Combine(healthCheckApi.LogInfo())...)
		// Report status to HAproxy
		err = c.AsyncWrite([]byte(status.String()), nil)
		if err != nil {
			logger.Error("Write", "message", err)
		}
	}()

	return gnet.None
}

func (s *Server) Close() error {
	return s.eng.Stop(context.Background())
}

//	Run the server event loop.
//
// ---------------------------------------------------------
func (s *Server) Run() error {
	return gnet.Run(s, s.addr, gnet.WithMulticore(true), gnet.WithTicker(true))
}
