
package log

import (
	"os"
	stdlog "log"
)


// ---------------------------------------------------------
//  Constants
// ---------------------------------------------------------


// Default flags to use.
const defaultFlags = stdlog.Lmsgprefix | stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds

// Prefixes
const (
	InfoPrefix    = "\033[1;34mINFO\033[0m   : "
	NoticePrefix  = "\033[1;36mNOTICE\033[0m : "
	WarningPrefix = "\033[1;33mWARN\033[0m   : "
	ErrorPrefix   = "\033[1;31mERROR\033[0m  : "
)


// ---------------------------------------------------------
//  Function wrappers
// ---------------------------------------------------------


// Define logWrapper prototype
//  Function that takes a format string and variadic number
//  of arguments (like printf)
// ---------------------------------------------------------
type logWrapper func(format string, args ...interface{})

// Create a wrapper function.
//  This creates a function wrapper around
//  stdlog.Logger.Printf()
// ---------------------------------------------------------
func createWrapper(logger *stdlog.Logger) logWrapper {
	return func(format string, args ...interface{}) {
		logger.Printf(format, args...)
	}
}

// Standard log wrapper
//  Wrapper around stdlog.Printf()
// ---------------------------------------------------------
func stdWrapper(format string, args ...interface{}) {
	stdlog.Printf(format, args...)
}


// ---------------------------------------------------------
//  Logger objects.
// ---------------------------------------------------------


var (
	// Info is standard logger. omitted here as we don't have direct access to the object.
	warningLogger 	*stdlog.Logger = stdlog.New(os.Stdout, WarningPrefix, defaultFlags)
	noticeLogger 	*stdlog.Logger = stdlog.New(os.Stdout, NoticePrefix, defaultFlags)
	errorLogger 	*stdlog.Logger = stdlog.New(os.Stdout, ErrorPrefix, defaultFlags)
)


// ---------------------------------------------------------
//  Initilize the module
// ---------------------------------------------------------


func init() {

	// Info is standard logger.
	// We are not allowed to access the standard
	// Logger object, so we have to use the standalone functions.

	stdlog.SetOutput(os.Stdout)
	stdlog.SetPrefix(InfoPrefix)
	stdlog.SetFlags(defaultFlags)
}
