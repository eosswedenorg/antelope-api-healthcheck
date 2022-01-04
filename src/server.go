package main

import (
    "fmt"
    "time"
    "strings"
    "strconv"
    "internal/utils"
    "github.com/eosswedenorg-go/eosapi"
    "github.com/eosswedenorg-go/haproxy"
    "github.com/eosswedenorg-go/tcp_server"
)

//  check_api - Validates head block time.
// ---------------------------------------------------------
func check_api(p eosapi.ReqParams, block_time float64) (haproxy.HealthCheckStatus, string) {

    info, err := eosapi.GetInfo(p)
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

    if diff > block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is lagging %.0f seconds", diff)
    } else if diff < -block_time {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because head block is %.0f seconds into the future", diff)
    }
    return haproxy.HealthCheckUp, "OK"
}

//  check_api_v2 (hyperion)
//    Validates block num diff between
//    nodeos and elasticsearch
// ---------------------------------------------------------
func check_api_v2(p eosapi.ReqParams, offset int64) (haproxy.HealthCheckStatus, string) {

    health, err := eosapi.GetHealth(p)
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
    if diff > offset {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Elastic is %d blocks behind", diff)
    } else if diff < -offset {
        return haproxy.HealthCheckDown,
            fmt.Sprintf("Taking offline because Elastic is %d blocks into the future", -1 * diff)
    }
    return haproxy.HealthCheckUp, "OK"
}

//  onTcpMessage callback function
// ---------------------------------------------------------

func onTcpMessage(c *tcp_server.Client, args string) {
    params := eosapi.ReqParams{}
    var block_time int = 10
    var version string = "v1"

    // Parse arguments.
    // -------------------
    split := strings.Split(strings.TrimSpace(args), "|")

    // 1. url (scheme + ip/domain + port)
    params.Url = split[0]

    // 2. Block time.
    if len(split) > 1 {
        p, err := strconv.ParseInt(split[1], 10, 32)
        if err == nil {
            block_time = int(p)
        }
    }

    // 3. Version
    if len(split) > 2 {
        version = split[2]
    }

    // 4. Host
    if len(split) > 3 {
        params.Host = split[3]
    }

    // Check api.
    // -------------------
    var status haproxy.HealthCheckStatus
    var msg string

    if version == "v2" {
        status, msg = check_api_v2(params, int64(block_time / 2))
    } else {
        version = "v1"
        status, msg = check_api(params, float64(block_time))
    }

    logger.Info("API Check", "version", version, "url", params.Url,
        "block", block_time / 2, "status", status)

    if status != haproxy.HealthCheckUp && len(msg) > 0 {
        logger.Warn("API Check Failed", "message", msg)
    }

    // Report status to HAproxy
    c.WriteString(fmt.Sprintln(status))
    c.Close()
}

//  spawnTcpServer
// ---------------------------------------------------------

func spawnTcpServer(addr string) {
    server := tcp_server.New(addr)
    server.OnMessage(onTcpMessage)
    server.Listen()
}
