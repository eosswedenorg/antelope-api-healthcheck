
package main

import (
	"os"
	"internal/pid"
	"internal/log"
	"github.com/pborman/getopt/v2"
)

//  Command line flags
// ---------------------------------------------------------

var logFile string
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

func openlog(file string) *os.File {

	fd, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err.Error())
	}
	return fd
}

//  main
// ---------------------------------------------------------
func main() {

	var version bool
	var logfd *os.File

	// Command line parsing
	getopt.FlagLong(&version, "version", 'v', "Print version")
	getopt.FlagLong(&logFile, "log", 'l', "Path to log file", "file")
	getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
	getopt.Parse()

	if version {
		print("Version: v1.0\n")
		return;
	}

	if len(logFile) > 0 {
		logfd = openlog(logFile)
		log.SetWriter(logfd)
	}

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
