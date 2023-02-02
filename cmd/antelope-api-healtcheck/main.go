package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eosswedenorg-go/pid"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/server"
	"github.com/eosswedenorg/antelope-api-healthcheck/internal/utils"
	log "github.com/inconshreveable/log15"
	"github.com/pborman/getopt/v2"
)

//  Command line flags
// ---------------------------------------------------------

var (
	logFile string
	pidFile string
)

//  Global variables
// ---------------------------------------------------------

// Version string, should be updated by the go linker (by passing "-X main.VersionString=value" to the linker)
// see: https://pkg.go.dev/cmd/link and
var VersionString string = "-"

// File descriptor to the current log file.
var logfd *os.File

var (
	logfmt log.Format
	logger log.Logger

	// TCP Server
	srv *server.Server
)

//	argv_listen_addr
//	  Parse listen address from command line.
//
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
	fd, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		logger.Error(err.Error())
	}

	// Try close if old descriptor is defined.
	if logfd != nil {
		if err = logfd.Close(); err != nil {
			logger.Error(err.Error())
		}
	}

	// Update variable and set log writer.
	logfd = fd
	logger.SetHandler(log.StreamHandler(logfd, logfmt))
}

//	signalEventLoop()
//	  Initialize event channel for OS signals
//	  and runs an event loop.
//
// ---------------------------------------------------------
func signalEventLoop() {
	// Setup a channel
	sig_ch := make(chan os.Signal, 1)

	// subscribe to SIGHUP signal.
	signal.Notify(sig_ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	// Event loop
	func() {
		var run bool = true
		for run {
			// Block until we get a signal.
			sig := <-sig_ch

			l := logger.New("signal", sig)

			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				l.Info("Program was asked to terminate.")
				run = false

				// Tell the server to close.
				err := srv.Close()
				if err != nil {
					l.Error("Failed to close server", "error", err)
				}
			// SIGHUP is sent when logfile is rotated.
			case syscall.SIGHUP:
				msg := "Logfile was rotated: "

				if logfd != nil {
					setLogFile()
					msg += "Filedescriptor was updated"
				} else {
					msg += "No Filedescriptor to update (most likely uses standard out/err streams)"
				}

				l.Info(msg)
			default:
				l.Warn("Unknown signal")
			}
		}
	}()
}

//	main
//
// ---------------------------------------------------------
func main() {
	var version bool
	var usage bool
	var logFormatter *string

	logger = log.Root()

	// Command line parsing
	getopt.SetParameters("[ip] [port]")
	getopt.FlagLong(&usage, "help", 'h', "Print this help text")
	getopt.FlagLong(&version, "version", 'v', "Print version")
	getopt.FlagLong(&logFile, "log", 'l', "Path to log file", "file")
	getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
	logFormatter = getopt.EnumLong("log-format", 0, []string{"term", "logfmt", "json", "json-pretty"}, "", "Log format to use: term,logfmt,json,json-pretty")

	getopt.Parse()

	if usage {
		getopt.Usage()
		return
	}

	if version {
		fmt.Printf("Version: %s\n", VersionString)
		return
	}

	logfmt = utils.ParseLogFormatter(*logFormatter)

	// Open logfile.
	if len(logFile) > 0 {
		setLogFile()
	} else {
		logger.SetHandler(log.StreamHandler(os.Stdout, logfmt))
	}

	logger.Info("Process is starting", "pid", pid.Get())

	if len(pidFile) > 0 {
		logger.Info("Writing pidfile", "file", pidFile)
		err := pid.Save(pidFile)
		if err != nil {
			logger.Error("Failed to write pidfile", "msg", err)
		}
	}

	// Create server
	srv = server.New(argv_listen_addr(), time.Second*10)

	// Run signal event loop in its own goroutine
	go signalEventLoop()

	// Run server
	if err := srv.Run(); err != nil {
		logger.Error("Server failed to shutdown", "message", err)
	}
}
