package server

import (
	"testing"

	"github.com/eosswedenorg/antelope-api-healthcheck/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestParseRequest_WithInvalidApi(t *testing.T) {
	api, err := ParseRequest("invalid|http://api.example.com")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid API 'invalid'")
	assert.Nil(t, api)
}

func TestParseRequest_WithInvalidParams(t *testing.T) {
	api, err := ParseRequest("v1")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid number of parameters in agent request")
	assert.Nil(t, api)
}

//  AntelopeV1
// --------------------------------

func TestParseRequest_AntelopeV1(t *testing.T) {
	expected := api.NewAntelopeV1("http://api.example.com", "", 5)

	api, err := ParseRequest("v1|http://api.example.com")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_AntelopeV1WithBlockNumber(t *testing.T) {
	expected := api.NewAntelopeV1("http://api.example.com", "", 1000)

	api, err := ParseRequest("v1|http://api.example.com|2000")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_AntelopeV1Full(t *testing.T) {
	expected := api.NewAntelopeV1("http://api.example.com", "http://host.example.com", 500)

	api, err := ParseRequest("v1|http://api.example.com|1000|http://host.example.com")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

//  AntelopeV2
// --------------------------------

func TestParseRequest_AntelopeV2(t *testing.T) {
	expected := api.NewAntelopeV2("http://api.v2.example.com", "", 10)

	api, err := ParseRequest("v2|http://api.v2.example.com")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_AntelopeV2WithOffset(t *testing.T) {
	expected := api.NewAntelopeV2("http://api.v2.example.com", "", 1000)

	api, err := ParseRequest("v2|http://api.v2.example.com|1000")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_AntelopeV2Full(t *testing.T) {
	expected := api.NewAntelopeV2("http://api.v2.example.com", "http://host.example.com", 1000)

	api, err := ParseRequest("v2|http://api.v2.example.com|1000|http://host.example.com")
	assert.NoError(t, err)

	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

// AtomicAsset
// --------------------------------

func TestParseRequest_AtomicAsset(t *testing.T) {
	expected := api.NewAtomicAsset("http://api.atomicassets.io", "", 5)

	api, err := ParseRequest("atomic|http://api.atomicassets.io")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_AtomicAssetWithBlockTime(t *testing.T) {
	expected := api.NewAtomicAsset("http://api.atomicassets.io", "", 256)

	api, err := ParseRequest("atomic|http://api.atomicassets.io|512")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseRequest_DebugApi(t *testing.T) {
	expected := api.NewDebugApi("some_api_call")

	api, err := ParseRequest("debug|some_api_call")
	assert.NoError(t, err)
	assert.Equal(t, expected.LogInfo(), api.LogInfo())
}
