
package api

import (
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

type ApiInterface interface {

    // Returns Logging information
    LogInfo() LogParams

    // Call api and validate it's status.
    Call() (agentcheck.Response, string)
}
