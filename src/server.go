package main

import (
    "fmt"
	"time"
	"strings"
	"strconv"
	"./log"
	"./haproxy"
	"./eosapi"
    "github.com/firstrow/tcp_server"
	"github.com/pborman/getopt/v2"
)

//  check_api - Validates head block time.
// ---------------------------------------------------------
func check_api(url string, block_time float64) (haproxy.HealthCheckStatus, string) {

	info, err := eosapi.GetInfo(url)
	if err != nil {
		msg := fmt.Sprintf("%s", err);
		return haproxy.HealthCheckFailed, msg
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
func check_api_v2(url string, offset int64) (haproxy.HealthCheckStatus, string) {

	health, err := eosapi.GetHealth(url)
	if err != nil {
		msg := fmt.Sprintf("%s", err);
		return haproxy.HealthCheckFailed, msg
	}

	// Fetch elasticsearch and nodeos block numbers from json.
	var es_block int64 = 0
	var node_block int64 = 0

	for _, v := range health.Health {
		if v.Name == "Elasticsearch" {
			es_block = (int64) (v.Data["last_indexed_block"].(float64))
		} else if v.Name == "NodeosRPC" {
			node_block = (int64) (v.Data["head_block_num"].(float64))
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

//  Command line flags
// ---------------------------------------------------------

var pidFile string

//  argv_listen_addr
//    Parse listen address from command line.
// ---------------------------------------------------------
func argv_listen_addr() string {

	var addr string

	argv := getopt.Args()
	if len(argv) > 0 {
		addr = argv[0]
	} else {
		addr = "127.0.0.1"
	}

	addr += ":"
	if len(argv) > 1 {
		addr += argv[1]
	} else {
		addr += "1337"
	}

	return addr
}

//  main
// ---------------------------------------------------------
func main() {

	// Command line parsing
	getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
	getopt.Parse()

	if len(pidFile) > 0 {
		log.Info("Writing pidfile: %s", pidFile)
	}

    server := tcp_server.New(argv_listen_addr())

	// TCP Client connect.
    server.OnNewClient(func(c *tcp_server.Client) {
        //fmt.Println("# Client connected")
    });

	// TCP Client sends message.
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		var url string
		var block_time int = 10
		var version string = "v1"

		// Parse host + port.
		split := strings.Split(strings.TrimSpace(message), "|")

		url = split[0]
		if len(split) > 1 {
			p, err := strconv.ParseInt(split[1], 10, 32)
			if err == nil {
				block_time = int(p)
			}
		}
		if len(split) > 2 {
			version = split[2]
		}

		// Check api.
		var status haproxy.HealthCheckStatus
		var msg string

		if version == "v2" {
			status, msg = check_api_v2(url, int64(block_time / 2))
		} else {
			version = "v1"
			status, msg = check_api(url, float64(block_time))
		}

		log.Info("Status %s - %s (%d blocks): %s",
			 version, url, block_time / 2, status)

		if status != haproxy.HealthCheckUp && len(msg) > 0 {
			log.Warning(msg)
		}

		// Report status to HAproxy
		c.Send(fmt.Sprintln(status))
        c.Close()
	});

	// TCP Client disconnect.
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		//fmt.Println("# Client disconnected")
	});

    server.Listen()
}
