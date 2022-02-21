
package api

import (
    "fmt"
    "time"
    "github.com/eosswedenorg-go/haproxy"
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

func (e EosioV1) Name() string {
    return "v1"
}

func (e EosioV1) Call() (haproxy.HealthCheckStatus, string) {

    info, err := eosapi.GetInfo(e.params)
    if err != nil {
        msg := fmt.Sprintf("%s", err);
        return haproxy.HealthCheckFailed, msg
    }

    // Check HTTP Status Code
    if info.HTTPStatusCode > 299 {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because %v was received from backend", info.HTTPStatusCode)
    }

    // Validate head block.
    now  := time.Now().In(time.UTC)
    diff := now.Sub(info.HeadBlockTime).Seconds()

    if diff > e.block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is lagging %.0f seconds", diff)
    } else if diff < -e.block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is %.0f seconds into the future", diff)
    }
    return haproxy.HealthCheckUp, "OK"
}
