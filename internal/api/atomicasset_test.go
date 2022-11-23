package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eosswedenorg-go/haproxy/agentcheck"
	"github.com/stretchr/testify/assert"
)

func TestAtomicAssetFactory(t *testing.T) {
	api := AtomicAssetFactory(ApiArguments{
		Url:       "https://atomic.example.com",
		NumBlocks: 120,
	})

	expected := NewAtomicAsset("https://atomic.example.com", 60)

	assert.IsType(t, expected, api)
	assert.Equal(t, expected.client.URL, api.(AtomicAsset).client.URL)
	assert.Equal(t, expected.client.Host, api.(AtomicAsset).client.Host)
	assert.Equal(t, expected.block_time, api.(AtomicAsset).block_time)
}

func TestAtomicAssetLogInfo(t *testing.T) {
	api := NewAtomicAsset("https://atomic.example.com", 120)

	expected := LogParams{"type", "atomicasset", "url", "https://atomic.example.com", "block_time", float64(120)}

	assert.Equal(t, expected, api.LogInfo())
}

func TestAtomicAssetSetTime(t *testing.T) {
	expected := time.Date(2019, 3, 18, 20, 29, 32, 0, time.UTC)

	api := NewAtomicAsset("", 60)
	// Assert that time is NOW (+-10 seconds)
	assert.InDelta(t, api.GetTime().Unix(), time.Now().In(time.UTC).Unix(), float64(10))

	api.SetTime(expected)
	assert.Equal(t, expected, api.GetTime())
}

func TestAtomicAssetJsonFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(`!//{invalid-json}!##`))
	}))

	api := NewAtomicAsset(srv.URL, 120)
	check, _ := api.Call()

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
	assert.Equal(t, expected, check)
}

func TestAtomicAssetHTTP500Down(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-type", "application/json; charset=utf-8")
		res.WriteHeader(500)
		res.Write([]byte(`{}`))
	}))

	api := NewAtomicAsset(srv.URL, 120)
	check, status := api.Call()

	assert.Equal(t, "Taking offline because 500 was received from backend", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}

func TestAtomicAssetLaggingUp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK"
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":2173612361,
                        "head_time":1759953927000
                    }
                },
                "query_time":1759953929542
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2025, 10, 8, 20, 7, 27, 0, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "OK", status)

	expected := agentcheck.NewStatusResponse(agentcheck.Up)
	assert.Equal(t, expected, check)
}

func TestAtomicAssetLaggingDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK"
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":213671263812,
                        "head_time":1533451894000
                    }
                },
                "query_time":1533451895542
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2018, 8, 5, 6, 53, 35, 0, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "Taking offline because head block is lagging 121 seconds", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}

func TestAtomicAssetInFutureUp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK"
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":213671263812,
                        "head_time":1728954676500
                    }
                },
                "query_time":1728954678231
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2024, 10, 15, 1, 9, 16, 500, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "OK", status)

	expected := agentcheck.NewStatusResponse(agentcheck.Up)
	assert.Equal(t, expected, check)
}

func TestAtomicAssetInFutureDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK"
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":213671263812,
                        "head_time":1041122824500
                    }
                },
                "query_time":1041122832231
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2002, 12, 29, 0, 45, 0o3, 500, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "Taking offline because head block is -121 seconds into the future", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}

func TestAtomicAssetRedisDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"OK"
                    },
                    "redis":{
                        "status":"DOWN"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":213671263812,
                        "head_time":1426072770500
                    }
                },
                "query_time":1426072775872
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2015, 3, 11, 11, 19, 30, 500, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "Taking offline because Redis reported 'DOWN'", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}

func TestAtomicAssetPostgresDown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/health" {
			payload := `{
                "success":true,
                "data":{
                    "version":"1.0.0",
                    "postgres":{
                        "status":"DOWN"
                    },
                    "redis":{
                        "status":"OK"
                    },
                    "chain":{
                        "status":"OK",
                        "head_block":213671263812,
                        "head_time":1562868371500
                    }
                },
                "query_time":156286837143
            }`

			res.Header().Add("Content-type", "application/json; charset=utf-8")
			res.Write([]byte(payload))
		}
	}))

	api := NewAtomicAsset(srv.URL, 120)
	api.SetTime(time.Date(2019, 7, 11, 18, 6, 11, 500, time.UTC))

	check, status := api.Call()

	assert.Equal(t, "Taking offline because Postgres reported 'DOWN'", status)

	expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
	assert.Equal(t, expected, check)
}
