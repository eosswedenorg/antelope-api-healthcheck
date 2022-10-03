
package api

import (
    "strings"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

type DebugApi struct {
    response agentcheck.Response
}

func parseResponse(resp string) (agentcheck.Response, error) {

    parts := strings.SplitN(resp, "#", 2)

    // Status with message
    if len(parts) > 1 {
        rtype := agentcheck.StatusMessageResponseType(parts[0])
        return agentcheck.NewStatusMessageResponse(rtype, parts[1]), nil
    }

    // Only status.
    rtype := agentcheck.StatusResponseType(parts[0])
    return agentcheck.NewStatusResponse(rtype), nil
}

func NewDebugApi(response string) DebugApi {

    resp, _ := parseResponse(response)

    return DebugApi{
        response: resp,
    }
}

func (d DebugApi) LogInfo() LogParams {
    return LogParams{
        "type", "Debug",
        "response", strings.TrimSpace(d.response.String()),
    }
}

func (d DebugApi) Call() (agentcheck.Response, string) {
    return d.response, ""
}
