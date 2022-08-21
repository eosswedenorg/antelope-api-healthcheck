package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    log "github.com/inconshreveable/log15"
    "github.com/eosswedenorg-go/pid"
    "github.com/pborman/getopt/v2"
)

//  Command line flags
// ---------------------------------------------------------

var logFile string
var pidFile string

//  Global variables
// ---------------------------------------------------------

// Version string, should be updated by the go linker (by passing "-X main.VersionString=value" to the linker)
// see: https://pkg.go.dev/cmd/link and
var VersionString string = "-"

// File descriptor to the current log file.
var logfd *os.File

var logger log.Logger

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
    logger.SetHandler(log.StreamHandler(logfd, log.LogfmtFormat()))
}

//  signalEventLoop()
//    Initialize event channel for OS signals
//    and runs an event loop.
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
            sig := <- sig_ch

            l := logger.New("signal", sig)

            switch sig {
            case syscall.SIGINT, syscall.SIGTERM :
                l.Info("Program was asked to terminate.")
                run = false
            // SIGHUP is sent when logfile is rotated.
            case syscall.SIGHUP :
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

//  main
// ---------------------------------------------------------
func main() {

    var version bool
    var usage bool
    var addr string;

    logger = log.New()

    // Command line parsing
    getopt.SetParameters("[ip] [port]")
    getopt.FlagLong(&usage, "help", 'h', "Print this help text")
    getopt.FlagLong(&version, "version", 'v', "Print version")
    getopt.FlagLong(&logFile, "log", 'l', "Path to log file", "file")
    getopt.FlagLong(&pidFile, "pid", 'p', "Path to pid file", "file")
    getopt.Parse()

    if usage {
        getopt.Usage()
        return
    }

    if version {
        fmt.Printf("Version: %s\n", VersionString)
        return;
    }

    // Open logfile.
    if len(logFile) > 0 {
        setLogFile()
    }

    logger.Info("Process is starting", "pid", pid.Get())

    if len(pidFile) > 0 {
        logger.Info("Writing pidfile", "file", pidFile)
        err := pid.Save(pidFile)
        if err != nil {
            logger.Error("Failed to write pidfile", "msg", err)
        }
    }

    addr = argv_listen_addr()

    // Start listening to TCP Connections
    err := spawnTcpServer(addr)
    if err == nil {
        logger.Info("TCP Server started", "addr", addr)

        // Run the signal event loop.
        signalEventLoop()
    } else {
        log.Error("Failed to start tcp server", "error", err)
    }

    logger.Info("Shutdown")
}
