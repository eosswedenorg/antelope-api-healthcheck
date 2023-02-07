package api

import (
	"context"
	"fmt"

	"github.com/eosswedenorg-go/atomicasset"
	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/utils"
)

type AtomicAsset struct {
	utils.Time

	url        string
	host       string
	block_time float64
}

func AtomicAssetFactory(args ApiArguments) ApiInterface {
	return NewAtomicAsset(args.Url, args.Host, float64(args.NumBlocks/2))
}

func NewAtomicAsset(url string, host string, block_time float64) AtomicAsset {
	return AtomicAsset{
		url:        url,
		host:       host,
		block_time: block_time,
	}
}

func (e AtomicAsset) LogInfo() LogParams {
	return LogParams{
		"type", "atomicasset",
		"url", e.url,
		"block_time", e.block_time,
	}
}

func (e AtomicAsset) Call(ctx context.Context) (agentcheck.Response, string) {
	client := atomicasset.NewWithContext(e.url, ctx)
	client.Host = e.host

	h, err := client.GetHealth()
	if err != nil {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")
		return resp, err.Error()
	}

	// Check HTTP Status Code
	if h.HTTPStatusCode > 299 {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		msg := "Taking offline because %v was received from backend"
		return resp, fmt.Sprintf(msg, h.HTTPStatusCode)
	}

	// Check postgres
	if h.Data.Postgres.Status != "OK" {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		msg := "Taking offline because Postgres reported '%s'"
		return resp, fmt.Sprintf(msg, h.Data.Postgres.Status)
	}

	// Check redis
	if h.Data.Redis.Status != "OK" {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		msg := "Taking offline because Redis reported '%s'"
		return resp, fmt.Sprintf(msg, h.Data.Redis.Status)
	}

	// Validate head block.
	diff := e.GetTime().Sub(h.Data.Chain.HeadTime.Time()).Seconds()

	if diff > e.block_time {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		msg := "Taking offline because head block is lagging %.0f seconds"
		return resp, fmt.Sprintf(msg, diff)
	} else if diff < -e.block_time {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		msg := "Taking offline because head block is %.0f seconds into the future"
		return resp, fmt.Sprintf(msg, diff)
	}

	return agentcheck.NewStatusResponse(agentcheck.Up), "OK"
}
