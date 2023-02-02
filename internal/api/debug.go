package api

import (
	"strings"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
)

type DebugApi struct {
	response agentcheck.Response
}

func parseResponse(resp string) agentcheck.Response {
	parts := strings.SplitN(resp, "#", 2)

	// Status with message
	if len(parts) > 1 {
		rtype := agentcheck.StatusMessageResponseType(parts[0])
		return agentcheck.NewStatusMessageResponse(rtype, parts[1])
	}

	// Only status.
	rtype := agentcheck.StatusResponseType(resp)
	return agentcheck.NewStatusResponse(rtype)
}

func DebugApiFactory(args ApiArguments) ApiInterface {
	return NewDebugApi(args.Url)
}

func NewDebugApi(response string) DebugApi {
	return DebugApi{
		response: parseResponse(response),
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
