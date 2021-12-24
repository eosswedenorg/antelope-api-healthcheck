
package main

import (
	"os"
	"os/signal"
	"syscall"
	"./log"
	"./pid"
	"github.com/pborman/getopt/v2"
)

//  Command line flags
// ---------------------------------------------------------

var logFile string
var pidFile string

//  Global variables
// ---------------------------------------------------------

// File descriptor to the current log file.
var logfd *os.File

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

func setLogFile() {

	// Open file
	fd, err := os.OpenFile(logFile, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Error(err.Error())
	}

	// Try close if old descriptor is defined.
	if logfd != nil {
		if err = logfd.Close(); err != nil {
			log.Error(err.Error())
		}
	}

	// Update variable and set log writer.
	logfd = fd
	log.SetWriter(logfd)
}

//  signalEventLoop()
//    Initialize event channel for OS signals
//    and runs an event loop in a separate thread.
// ---------------------------------------------------------
func signalEventLoop() {

	// Setup a channel
	sig_ch := make(chan os.Signal, 1)

	// subscribe to USR1 signal.
	signal.Notify(sig_ch, syscall.SIGUSR1)

	// Event loop (runs in a seperate thread)
	go func() {
		for {
			// Block until we get a signal.
			sig := <- sig_ch

			switch sig {
			// USR1 is sent when logfile is rotated.
			case syscall.SIGUSR1 :
				msg := "SIGUSR1 (Logfile was rotated): "

				if logfd != nil {
					setLogFile()
					msg += "Filedescriptor was updated"
				} else {
					msg += "No Filedescriptor to update (most likely uses standard out/err streams)"
				}

				log.Info(msg)
			default:
				log.Warning("Unknown signal %s", sig)
			}
		}
	}()
}

//  main
// ---------------------------------------------------------
func main() {

	// Command line parsing
	getopt.FlagLong(&logFile, "log", 'l', "Path to log file", "file")
	getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
	getopt.Parse()

	// Open logfile.
	if len(logFile) > 0 {
		setLogFile()
	}

	log.Info("Process is starting with PID: %d", pid.Get())

	if len(pidFile) > 0 {
		log.Info("Writing pidfile: %s", pidFile)
		_, err := pid.Save(pidFile)
		if err != nil {
			log.Error("Failed to write pidfile: %v", err)
		}
	}

	// Run the signal event loop.
	signalEventLoop()

	// Start listening to TCP Connections
	spawnTcpServer(argv_listen_addr());
}
