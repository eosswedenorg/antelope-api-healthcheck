
package log

import (
	"io"
	stdlog "log"
)

// Create and export our different log functions.
// ---------------------------------------------------------
var (
	Info logWrapper 	= stdWrapper
	Notice logWrapper 	= createWrapper(noticeLogger)
	Warning logWrapper  = createWrapper(warningLogger)
	Error logWrapper 	= createWrapper(errorLogger)
)

// SetWriter
// 	Configure the logger to use this writer.
// ---------------------------------------------------------
func SetWriter(handle io.Writer) {

	// Info is standard logger.
	stdlog.SetOutput(handle)

	noticeLogger.SetOutput(handle)
	warningLogger.SetOutput(handle)
	errorLogger.SetOutput(handle)
}
