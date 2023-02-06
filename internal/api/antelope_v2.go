package api

import (
	"context"
	"fmt"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/eosswedenorg-go/leapapi"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/utils"
)

type AntelopeV2 struct {
	client leapapi.Client
	offset int64
}

func AntelopeV2Factory(args ApiArguments) ApiInterface {
	return NewAntelopeV2(args.Url, args.Host, int64(args.NumBlocks))
}

func NewAntelopeV2(url string, host string, offset int64) AntelopeV2 {
	api := AntelopeV2{
		client: *leapapi.New(url),
		offset: offset,
	}

	api.client.Host = host

	return api
}

func (e AntelopeV2) LogInfo() LogParams {
	p := LogParams{
		"type", "antelope-v2",
		"url", e.client.Url,
	}

	if len(e.client.Host) > 0 {
		p.Add("host", e.client.Host)
	}

	p.Add("offset", e.offset)

	return p
}

func (e AntelopeV2) Call(ctx context.Context) (agentcheck.Response, string) {
	// TODO: Pass context
	health, err := e.client.GetHealth()
	if err != nil {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")
		return resp, err.Error()
	}

	// Fetch elasticsearch and nodeos block numbers from json.
	var es_block int64 = 0
	var node_block int64 = 0

	for _, v := range health.Health {
		if v.Name == "Elasticsearch" {
			es_block = utils.JsonGetInt64(v.Data["last_indexed_block"])
		} else if v.Name == "NodeosRPC" {
			node_block = utils.JsonGetInt64(v.Data["head_block_num"])
		}
	}

	// Error out if ether or both are zero.
	if es_block == 0 || node_block == 0 {
		msg := fmt.Sprintf("Failed to get Elasticsearch and/or nodeos "+
			"block numbers (es: %d, eos: %d)", es_block, node_block)

		resp := agentcheck.NewStatusMessageResponse(agentcheck.Fail, "")
		return resp, msg
	}

	// Check if ES is behind or in the future.
	diff := node_block - es_block
	if diff > e.offset {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		return resp, fmt.Sprintf("Taking offline because Elastic is %d blocks behind", diff)
	} else if diff < -e.offset {
		resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
		return resp, fmt.Sprintf("Taking offline because Elastic is %d blocks into the future", -1*diff)
	}
	return agentcheck.NewStatusResponse(agentcheck.Up), "OK"
}
