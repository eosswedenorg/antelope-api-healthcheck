
package api

import (
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

/**
 * Generic struct that is passed to factory functions
 * to configure the API request.
 */
type ApiArguments struct {
    Url string
    Host string
    NumBlocks int
}

type ApiInterface interface {

    // Returns Logging information
    LogInfo() LogParams

    // Call api and validate it's status.
    Call() (agentcheck.Response, string)
}
