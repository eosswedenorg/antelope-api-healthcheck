package api

import (
	"fmt"

	"github.com/eosswedenorg-go/eosapi"
	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg/eosio-api-healthcheck/internal/utils"
)

type EosioV1 struct {
	utils.Time
	client     eosapi.Client
	block_time float64
}

func EosioV1Factory(args ApiArguments) ApiInterface {
	return NewEosioV1(args.Url, args.Host, float64(args.NumBlocks/2))
}

func NewEosioV1(url string, host string, block_time float64) EosioV1 {
	api := EosioV1{
		client:     *eosapi.New(url),
		block_time: block_time,
	}

	api.client.Host = host

	return api
}

func (e EosioV1) LogInfo() LogParams {
	p := LogParams{
		"type", "eosio-v1",
		"url", e.client.Url,
	}

	if len(e.client.Host) > 0 {
		p.Add("host", e.client.Host)
	}

	p.Add("block_time", e.block_time)

	return p
}

func (e EosioV1) Call() (agentcheck.Response, string) {
	info, err := e.client.GetInfo()
	if err != nil {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
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
