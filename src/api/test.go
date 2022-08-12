
package api

import (
    "strings"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

type TestApi struct {
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

func NewTestApi(response string) TestApi {

    resp, _ := parseResponse(response)

    return TestApi{
        response: resp,
    }
}

func (t TestApi) LogInfo() LogParams {
    return LogParams{
        "type", "TestApi",
        "response", t.response,
    }
}

func (t TestApi) Call() (agentcheck.Response, string) {
    return t.response, ""
}
