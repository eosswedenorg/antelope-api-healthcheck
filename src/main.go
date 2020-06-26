
package main

import (
	"./log"
	"./pid"
	"github.com/pborman/getopt/v2"
)

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

	log.Info("Process is starting with PID: %d", pid.Get())

	if len(pidFile) > 0 {
		log.Info("Writing pidfile: %s", pidFile)
		_, err := pid.Save(pidFile)
		if err != nil {
			log.Error("Failed to write pidfile: %v", err)
		}
	}

	spawnTcpServer(argv_listen_addr());
}
