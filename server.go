package main

import (
	"os"
    "fmt"
	"time"
	"strings"
	"strconv"
	"./haproxy"
	"./eosapi"
    "github.com/firstrow/tcp_server"
)

//  check_api - Validates head block time.
// ---------------------------------------------------------
func check_api(host string, port int) (haproxy.HealthCheckStatus, string) {

	info, err := eosapi.GetInfo(host, port)
	if err != nil {
		msg := fmt.Sprintf("%s", err);
		return haproxy.HealthCheckFailed, msg
	}

	// Validate head block.
	now  := time.Now().In(time.UTC)
	diff := now.Sub(info.HeadBlockTime).Seconds()

	if diff > 10.0 {
		return haproxy.HealthCheckDown,
			fmt.Sprintf("Taking offline because head block is lagging %.0f seconds", diff)
	} else if diff < -10.0 {
		return haproxy.HealthCheckDown,
			fmt.Sprintf("Taking offline because head block is %.0f seconds into the future", diff)
	}
	return haproxy.HealthCheckUp, "OK"
}

//  argv_listen_addr
//    Parse listen address from command line.
// ---------------------------------------------------------
func argv_listen_addr() string {

	var addr string

	argv := os.Args[1:]
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

    server := tcp_server.New(argv_listen_addr())

	// TCP Client connect.
    server.OnNewClient(func(c *tcp_server.Client) {
        fmt.Println("# Client connected")
    });

	// TCP Client sends message.
	server.OnNewMessage(func(c *tcp_server.Client, message string) {
		var host string
		var port int = 80

		// Parse host + port.
		split := strings.Split(strings.TrimSpace(message), ":")

		host = split[0]
		if len(split) > 1 {
			p, err := strconv.ParseInt(split[1], 10, 32)
			if err == nil {
				port = int(p)
			}
		}

		// Check api.
		status, msg := check_api(host, port)

		fmt.Printf("API HealthCheck: %s, %s\n", status, msg)

		// Report status to HAproxy
		c.Send(fmt.Sprintln(status))
        c.Close()
	});

	// TCP Client disconnect.
	server.OnClientConnectionClosed(func(c *tcp_server.Client, err error) {
		fmt.Println("# Client disconnected")
	});

    server.Listen()
}
