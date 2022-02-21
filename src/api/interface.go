
package api

import (
    "github.com/eosswedenorg-go/haproxy"
)

type ApiInterface interface {

    // Returns Logging information
    LogInfo() LogParams

    // Call api and validate it's status.
    Call() (haproxy.HealthCheckStatus, string)
}
