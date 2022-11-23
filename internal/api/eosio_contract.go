package api

import (
	"fmt"

	"github.com/eosswedenorg-go/atomicasset"
	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg/eosio-api-healthcheck/internal/utils"
)

type EosioContract struct {
	utils.Time
	client     atomicasset.Client
	block_time float64
}

func EosioContractFactory(args ApiArguments) ApiInterface {
	return NewEosioContract(args.Url, float64(args.NumBlocks/2))
}

func NewEosioContract(url string, block_time float64) EosioContract {
	return EosioContract{
		client: atomicasset.Client{
			URL: url,
		},
		block_time: block_time,
	}
}

func (e EosioContract) LogInfo() LogParams {
	return LogParams{
		"type", "eosio-contract",
		"url", e.client.URL,
		"block_time", e.block_time,
	}
}

func (e EosioContract) Call() (agentcheck.Response, string) {
	h, err := e.client.GetHealth()
	if err != nil {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
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
