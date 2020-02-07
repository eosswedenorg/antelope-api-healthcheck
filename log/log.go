
package log

import "fmt"

// Colors
// ---------------------------------------------------------
const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

// Define LogFunc prototype
//  Function that takes a format string and variadic number
//  of arguments (like printf)
// ---------------------------------------------------------
type LogFunc func(format string, args ...interface{})

// Create a log function.
//  This is the base logging function. by providing a prefix
//  a new log function of LogFunc type will be created
//  appending "[ <prefix> ]" before the message.
// ---------------------------------------------------------
func logfn(prefix string) LogFunc {
	return func(format string, args ...interface{}) {
		format = "[" + prefix + "] " + format + "\n"
		fmt.Printf(format, args...)
	}
}

// Declare our different log functions.
// ---------------------------------------------------------
var Info LogFunc
var Notice LogFunc
var Warning LogFunc
var Error LogFunc

// Initilize log module
// ---------------------------------------------------------
func init() {

	// Initilize functions.
	Info 	= logfn(fmt.Sprintf(InfoColor, "INFO"))
	Notice 	= logfn(fmt.Sprintf(NoticeColor, "NOTICE"))
	Warning = logfn(fmt.Sprintf(WarningColor, "WARN"))
	Error 	= logfn(fmt.Sprintf(ErrorColor, "ERROR"))
}
