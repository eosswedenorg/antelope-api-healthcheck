
package api

import (
    "fmt"
    "time"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
    "github.com/eosswedenorg-go/eosapi"
)

type EosioV1 struct {
    params eosapi.ReqParams
    block_time float64
}

func NewEosioV1(params eosapi.ReqParams, block_time float64) EosioV1 {
    return EosioV1{
        params: params,
        block_time: block_time,
    }
}

func (e EosioV1) LogInfo() LogParams {
    p := LogParams{
        "type", "eosio-v1",
        "url", e.params.Url,
    }

    if len(e.params.Host) > 0 {
        p.Add("host", e.params.Host)
    }

    p.Add("block_time", e.block_time)

    return p
}

func (e EosioV1) Call() (agentcheck.Response, string) {

    info, err := eosapi.GetInfo(e.params)
    if err != nil {
        resp := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
        return resp, err.Error()
    }

    // Check HTTP Status Code
    if info.HTTPStatusCode > 299 {

        resp := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")

        msg := "Taking offline because %v was received from backend"
        return resp, fmt.Sprintf(msg, info.HTTPStatusCode)
    }

    // Validate head block.
    now  := time.Now().In(time.UTC)
    diff := now.Sub(info.HeadBlockTime).Seconds()

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
