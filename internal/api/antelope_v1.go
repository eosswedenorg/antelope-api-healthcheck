package api

import (
	"fmt"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg-go/leapapi"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/utils"
)

type AntelopeV1 struct {
	utils.Time
	client     leapapi.Client
	block_time float64
}

func AntelopeV1Factory(args ApiArguments) ApiInterface {
	return NewAntelopeV1(args.Url, args.Host, float64(args.NumBlocks/2))
}

func NewAntelopeV1(url string, host string, block_time float64) AntelopeV1 {
	api := AntelopeV1{
		client:     *leapapi.New(url),
		block_time: block_time,
	}

	api.client.Host = host

	return api
}

func (e AntelopeV1) LogInfo() LogParams {
	p := LogParams{
		"type", "antelope-v1",
		"url", e.client.Url,
	}

	if len(e.client.Host) > 0 {
		p.Add("host", e.client.Host)
	}

	p.Add("block_time", e.block_time)

	return p
}

func (e AntelopeV1) Call() (agentcheck.Response, string) {
	info, err := e.client.GetInfo()
	if err != nil {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")
		return resp, err.Error()
	}

	// Validate head block.
	diff := e.GetTime().Sub(info.HeadBlockTime).Seconds()

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
