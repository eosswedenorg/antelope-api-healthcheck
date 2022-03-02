
package api

import (
    "fmt"
    "time"
    "github.com/eosswedenorg-go/haproxy"
    contract_api "github.com/eosswedenorg-go/eos-contract-api-client"
)

type EosioContract struct {
    client contract_api.Client
    block_time float64
}

func NewEosioContract(url string, block_time float64) EosioContract {
    return EosioContract{
        client: contract_api.Client{
            Url: url,
        },
        block_time: block_time,
    }
}

func (e EosioContract) LogInfo() LogParams {
    return LogParams{
        "type", "eosio-contract",
        "url", e.client.Url,
        "block_time", e.block_time,
    }
}

//  check_api - Validates head block time.
// ---------------------------------------------------------
func (e EosioContract) Call() (haproxy.HealthCheckStatus, string) {

    h, err := e.client.GetHealth()
    if err != nil {
        msg := fmt.Sprintf("%s", err);
        return haproxy.HealthCheckFailed, msg
    }

    // Check HTTP Status Code
    if h.HTTPStatusCode > 299 {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because %v was received from backend", h.HTTPStatusCode)
    }

    // Check postgres
    if h.Data.Postgres.Status != "OK" {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Postgres reported '%s'", h.Data.Postgres.Status)
    }

    // Check redis
    if h.Data.Redis.Status != "OK" {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Redis reported '%s'", h.Data.Redis.Status)
    }

    // Validate head block.
    now  := time.Now().In(time.UTC)
    diff := now.Sub(h.Data.Chain.HeadTime).Seconds()

    if diff > e.block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is lagging %.0f seconds", diff)
    } else if diff < -e.block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is %.0f seconds into the future", diff)
    }

    return haproxy.HealthCheckUp, "OK"
}
