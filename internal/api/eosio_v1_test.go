package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/stretchr/testify/assert"
)

func TestAntelopeV1_Factory(t *testing.T) {
	api := AntelopeV1Factory(ApiArguments{
		Url:       "https://api.v1.example.com",
		Host:      "host.example.com",
		NumBlocks: 120,
	})

	expected := NewAntelopeV1("https://api.v1.example.com", "host.example.com", 60)

	assert.IsType(t, expected, api)
	assert.Equal(t, expected.client.Url, api.(AntelopeV1).client.Url)
	assert.Equal(t, expected.client.Host, api.(AntelopeV1).client.Host)
	assert.Equal(t, expected.block_time, api.(AntelopeV1).block_time)
}

func TestAntelopeV1_LogInfo(t *testing.T) {
	api := NewAntelopeV1("https://api.v1.example.com", "host.example.com", 120)

	expected := LogParams{"type", "antelope-v1", "url", "https://api.v1.example.com", "host", "host.example.com", "block_time", float64(120)}

	assert.Equal(t, expected, api.LogInfo())
}

func TestAntelopeV1_SetTime(t *testing.T) {
	expected := time.Date(2022, 2, 24, 13, 38, 0, 0, time.UTC)

	api := NewAntelopeV1("", "", 60)
	// Assert that time is NOW (+-10 seconds)
	assert.InDelta(t, api.GetTime().Unix(), time.Now().In(time.UTC).Unix(), float64(10))

	api.SetTime(expected)
	assert.Equal(t, expected, api.GetTime())
}

func TestAntelopeV1_JsonFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, err := res.Write([]byte(`!//{invalid-json}!##`))
		assert.NoError(t, err)
	}))

	api := NewAntelopeV1(srv.URL, "", 120)
	check, _ := api.Call()

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
	assert.Equal(t, expected, check)
}

func TestAntelopeV1_HTTP500Failed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		_, err := res.Write([]byte(`{}`))
		assert.NoError(t, err)
	}))

	api := NewAntelopeV1(srv.URL, "", 120)
	check, status := api.Call()

	assert.Equal(t, "server returned HTTP 500 Internal Server Error", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
	assert.Equal(t, expected, check)
}

func TestAntelopeV1_LaggingUp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/v1/chain/get_info" {
			info := `{
                "server_version": "8f613ec9",
                "head_block_num": 7272812,
                "head_block_time": "2022-02-24T13:37:00"
            }`

			_, err := res.Write([]byte(info))
			assert.NoError(t, err)
		}
	}))

	api := NewAntelopeV1(srv.URL, "", 60)
	api.SetTime(time.Date(2022, 2, 24, 13, 38, 0, 0, time.UTC))
	check, status := api.Call()

	assert.Equal(t, "OK", status)

	expected := agentcheck.NewStatusResponse(agentcheck.Up)
	assert.Equal(t, expected, check)
}

func TestAntelopeV1_LaggingDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/v1/chain/get_info" {
			info := `{
                "server_version": "9a607cce",
                "head_block_num": 87263,
                "head_block_time": "2018-01-01T13:37:01"
            }`

			_, err := res.Write([]byte(info))
			assert.NoError(t, err)
		}
	}))

	api := NewAntelopeV1(srv.URL, "", 60)
	api.SetTime(time.Date(2018, time.January, 1, 13, 38, 2, 0, time.UTC))
	check, status := api.Call()

	assert.Equal(t, "Taking offline because head block is lagging 61 seconds", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}

func TestAntelopeV1_TimeInFutureUP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/v1/chain/get_info" {
			info := `{
                "server_version": "d1bec8d3",
                "head_block_num": 548847,
                "head_block_time": "2020-09-22T09:32:00"
            }`

			_, err := res.Write([]byte(info))
			assert.NoError(t, err)
		}
	}))

	api := NewAntelopeV1(srv.URL, "", 120)
	api.SetTime(time.Date(2020, 9, 22, 9, 30, 0, 0, time.UTC))
	check, status := api.Call()

	assert.Equal(t, "OK", status)

	expected := agentcheck.NewStatusResponse(agentcheck.Up)
	assert.Equal(t, expected, check)
}

func TestAntelopeV1_TimeInFutureDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/v1/chain/get_info" {
			info := `{
                "server_version": "c879d231",
                "head_block_num": 2637621,
                "head_block_time": "2019-04-14T12:02:01"
            }`

			_, err := res.Write([]byte(info))
			assert.NoError(t, err)
		}
	}))

	api := NewAntelopeV1(srv.URL, "", 120)
	api.SetTime(time.Date(2019, time.April, 14, 12, 0, 0, 0, time.UTC))
	check, status := api.Call()

	assert.Equal(t, "Taking offline because head block is -121 seconds into the future", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}
