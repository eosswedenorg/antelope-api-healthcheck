
package haproxy;

// All supported health check values for HAproxy.
// See https://cbonte.github.io/haproxy-dconv/1.7/configuration.html#5.2-agent-check
type HealthCheckStatus string
const (
	HealthCheckUp 		= "up"
	HealthCheckDown 	= "down"
	HealthCheckMaint 	= "maint"
	HealthCheckReady	= "ready"
	HealthCheckDrain	= "drain"
	HealthCheckFailed 	= "failed"
	HealthCheckStopped 	= "Stopped"
)
