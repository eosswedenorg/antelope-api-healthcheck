
package main

import (
    log "github.com/inconshreveable/log15"
)

func parseLogFormatter(name string) log.Format {

    switch name {
    case "logfmt" :
        return log.LogfmtFormat()
    case "json" :
        return log.JsonFormat()
    case "json-pretty" :
        return log.JsonFormatEx(true, true)
    default :
        return log.TerminalFormat()
    }
}
