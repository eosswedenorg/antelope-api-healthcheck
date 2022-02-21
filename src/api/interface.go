
package api

import (
    "github.com/eosswedenorg-go/haproxy"
)

type ApiInterface interface {

    // Name of the api
    Name() string

    // Call api and validate it's status.
    Call() (haproxy.HealthCheckStatus, string)
}
