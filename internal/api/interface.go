package api

import (
	"context"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
)

/**
 * Generic struct that is passed to factory functions
 * to configure the API request.
 */
type ApiArguments struct {
	Url       string
	Host      string
	NumBlocks int
}

/**
 * Factory function
 *
 * Each API must implement this function and process `args`
 * returing a instance of it's implementation of the ApiInterface
 */
type Factory func(args ApiArguments) ApiInterface

type ApiInterface interface {
	// Returns Logging information
	LogInfo() LogParams

	// Call api and validate it's status.
	Call(ctx context.Context) (agentcheck.Response, string)
}
