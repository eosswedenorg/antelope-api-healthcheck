package server

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/api"
	log "github.com/inconshreveable/log15"
	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	gnet.BuiltinEventEngine

	eng gnet.Engine

	// Address to bind to.
	addr string

	// Number of connections between each OnTick()
	num_conn uint64

	// Time between each call to OnTick()
	tick_interval time.Duration
}

type Option func(*Server)

func New(addr string, options ...Option) *Server {
	s := &Server{
		addr: fmt.Sprintf("tcp://%s", addr),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func WithTick(interval time.Duration) Option {
	return func(s *Server) {
		s.tick_interval = interval
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

func (s *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	atomic.AddUint64(&s.num_conn, 1)
	return nil, gnet.None
}

//	OnClose callback function
//
// ---------------------------------------------------------
func (s *Server) OnClose(c gnet.Conn, err error) gnet.Action {
	if err != nil {
		log.Error("TCP Close", "error", err)
	}
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
	log.Info("Server info", log.Ctx{
		"connections":         atomic.LoadUint64(&s.num_conn),
		"current_connections": s.eng.CountConnections(),
	})
	atomic.StoreUint64(&s.num_conn, 0)
	return s.tick_interval, gnet.None
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
		// Make a context with 30 sec timeout per default. Should be "enough" for most cases.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		t := time.Now()
		status, msg := healthCheckApi.Call(ctx)
		req_time := time.Since(t)

		params := api.LogParams{}
		params.Add("status", strings.TrimSpace(status.String()))
		params.Add("duration", req_time)
		params.Add("duration_us", req_time.Microseconds())

		if msg != "OK" && len(msg) > 0 {
			params.Add("error", msg)
		}

		logger.Info("API Check", params.Combine(healthCheckApi.LogInfo())...)
		// Report status to HAproxy
		err = c.AsyncWrite([]byte(status.String()), func(c gnet.Conn, err error) error {
			return c.Close()
		})

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
