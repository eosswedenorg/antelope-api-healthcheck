
package api

import (
    "github.com/eosswedenorg-go/haproxy"
)

type ApiInterface interface {
    Name() string
    Call() (haproxy.HealthCheckStatus, string)
}
