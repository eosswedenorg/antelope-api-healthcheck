
package internal

import (
    // "fmt"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/eosswedenorg/eosio-api-healthcheck/internal/api"
)

func TestParseWithInvalidApi(t *testing.T) {

    api, err := ParseRequest("invalid|http://api.example.com")
    assert.Error(t, err)
    assert.Equal(t, err.Error(), "invalid API 'invalid'")
    assert.Nil(t, api)
}

func TestParseWithInvalidParams(t *testing.T) {

    api, err := ParseRequest("v1")
    assert.Error(t, err)
    assert.Equal(t, err.Error(), "invalid number of parameters in agent request")
    assert.Nil(t, api)
}

//  EosioV1
// --------------------------------

func TestParseEosioV1(t *testing.T) {

    expected := api.NewEosioV1("http://api.example.com", "", 5)

    api, err := ParseRequest("v1|http://api.example.com")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseEosioV1WithBlockNumber(t *testing.T) {

    expected := api.NewEosioV1("http://api.example.com", "", 1000)

    api, err := ParseRequest("v1|http://api.example.com|2000")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}


func TestParseEosioV1Full(t *testing.T) {

    expected := api.NewEosioV1("http://api.example.com", "http://host.example.com", 500)

    api, err := ParseRequest("v1|http://api.example.com|1000|http://host.example.com")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

//  EosioV2
// --------------------------------

func TestParseEosioV2(t *testing.T) {

    expected := api.NewEosioV2("http://api.v2.example.com", "", 10)

    api, err := ParseRequest("v2|http://api.v2.example.com")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseEosioV2WithOffset(t *testing.T) {

    expected := api.NewEosioV2("http://api.v2.example.com", "", 1000)

    api, err := ParseRequest("v2|http://api.v2.example.com|1000")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseEosioV2Full(t *testing.T) {

    expected := api.NewEosioV2("http://api.v2.example.com", "http://host.example.com", 1000)

    api, err := ParseRequest("v2|http://api.v2.example.com|1000|http://host.example.com")
    assert.NoError(t, err)

    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

// EosioContract
// --------------------------------

func TestParseEosioContract(t *testing.T) {

    expected := api.NewEosioContract("http://api.contract.example.com", 5)

    api, err := ParseRequest("contract|http://api.contract.example.com")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseEosioContractWithBlockTime(t *testing.T) {

    expected := api.NewEosioContract("http://api.contract.example.com", 256)

    api, err := ParseRequest("contract|http://api.contract.example.com|512")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}

func TestParseDebugApi(t *testing.T) {

    expected := api.NewDebugApi("some_api_call")

    api, err := ParseRequest("debug|some_api_call")
    assert.NoError(t, err)
    assert.Equal(t, expected.LogInfo(), api.LogInfo())
}
