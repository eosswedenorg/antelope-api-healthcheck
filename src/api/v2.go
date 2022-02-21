
package api

import (
    "fmt"
    "internal/utils"
    "github.com/eosswedenorg-go/haproxy"
    "github.com/eosswedenorg-go/eosapi"
)

type EosioV2 struct {
    params eosapi.ReqParams
    offset int64
}

func NewEosioV2(params eosapi.ReqParams, offset int64) EosioV2 {
    return EosioV2{
        params: params,
        offset: offset,
    }
}

func (e EosioV2) LogInfo() LogParams {
    p := LogParams{
        "type", "eosio-v2",
        "url", e.params.Url,
    }

    if len(e.params.Host) > 0 {
        p.Add("host", e.params.Host)
    }

    p.Add("offset", e.offset)

    return p
}

func (e EosioV2) Call() (haproxy.HealthCheckStatus, string) {

    health, err := eosapi.GetHealth(e.params)
    if err != nil {
        msg := fmt.Sprintf("%s", err);
        return haproxy.HealthCheckFailed, msg
    }

    // Check HTTP Status Code
    if health.HTTPStatusCode > 299 {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because %v was received from backend", health.HTTPStatusCode)
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
        msg := fmt.Sprintf("Failed to get Elasticsearch and/or nodeos " +
            "block numbers (es: %d, eos: %d)", es_block, node_block)
        return haproxy.HealthCheckFailed, msg
    }

    // Check if ES is behind or in the future.
    diff := node_block - es_block;
    if diff > e.offset {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Elastic is %d blocks behind", diff)
    } else if diff < -e.offset {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Elastic is %d blocks into the future", -1 * diff)
    }
    return haproxy.HealthCheckUp, "OK"
}
