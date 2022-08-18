
package api

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

func TestEosioV2LogInfo(t *testing.T) {

    api := NewEosioV2("https://api.v2.example.com", "host.example.com", 120)

    expected := LogParams{"type","eosio-v2","url","https://api.v2.example.com","host","host.example.com","offset",int64(120)}

    assert.Equal(t, expected, api.LogInfo())
}

func TestEosioV2JsonFailure(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.Write([]byte(`!//{invalid-json}!##`))
    }))

    api := NewEosioV2(srv.URL, "", 120)
    check, _ := api.Call()

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2HTTP500Down(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(500)
        res.Write([]byte(`{}`))
    }))

    api := NewEosioV2(srv.URL, "", 120)
    check, status := api.Call()

    assert.Equal(t, "Taking offline because 500 was received from backend", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2LaggingUp(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "OK",
                        "service_data": {
                            "head_block_num": 263148621,
                            "head_block_time": "2022-08-17T14:16:36.000",
                            "time_offset": 190,
                            "last_irreversible_block": 263148296,
                            "chain_id": "f8c74ccb7f9dea6f26a6d7f786809ddd1bce9fada3867f567dd83691b5348534"
                        },
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "OK",
                        "service_data": {
                            "last_indexed_block": 263148121,
                            "total_indexed_blocks": 263148121,
                            "active_shards": "100.0%"
                        },
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 500)
    check, status := api.Call()

    assert.Equal(t, "OK", status)

    expected := agentcheck.NewStatusResponse(agentcheck.Up)
    assert.Equal(t, expected, check)
}

func TestEosioV2LaggingDown(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "OK",
                        "service_data": {
                            "head_block_num": 263148621,
                            "head_block_time": "2022-08-17T14:16:36.000",
                            "time_offset": 190,
                            "last_irreversible_block": 263148296,
                            "chain_id": "f8c74ccb7f9dea6f26a6d7f786809ddd1bce9fada3867f567dd83691b5348534"
                        },
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "OK",
                        "service_data": {
                            "last_indexed_block": 263148121,
                            "total_indexed_blocks": 263148121,
                            "active_shards": "100.0%"
                        },
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 499)
    check, status := api.Call()

    assert.Equal(t, "Taking offline because Elastic is 500 blocks behind", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2LaggingESInFutureUP(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "OK",
                        "service_data": {
                            "head_block_num": 263148621,
                            "head_block_time": "2022-08-17T14:16:36.000",
                            "time_offset": 190,
                            "last_irreversible_block": 263148296,
                            "chain_id": "f8c74ccb7f9dea6f26a6d7f786809ddd1bce9fada3867f567dd83691b5348534"
                        },
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "OK",
                        "service_data": {
                            "last_indexed_block": 263148821,
                            "total_indexed_blocks": 263148821,
                            "active_shards": "100.0%"
                        },
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 200)
    check, status := api.Call()

    assert.Equal(t, "OK", status)

    expected := agentcheck.NewStatusResponse(agentcheck.Up)
    assert.Equal(t, expected, check)
}

func TestEosioV2LaggingESInFutureDown(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "OK",
                        "service_data": {
                            "head_block_num": 263148621,
                            "head_block_time": "2022-08-17T14:16:36.000",
                            "time_offset": 190,
                            "last_irreversible_block": 263148296,
                            "chain_id": "f8c74ccb7f9dea6f26a6d7f786809ddd1bce9fada3867f567dd83691b5348534"
                        },
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "OK",
                        "service_data": {
                            "last_indexed_block": 263148822,
                            "total_indexed_blocks": 263148822,
                            "active_shards": "100.0%"
                        },
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 200)
    check, status := api.Call()

    assert.Equal(t, "Taking offline because Elastic is 201 blocks into the future", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2ElasticsFailed(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "OK",
                        "service_data": {
                            "head_block_num": 263148621,
                            "head_block_time": "2022-08-17T14:16:36.000",
                            "time_offset": 190,
                            "last_irreversible_block": 263148296,
                            "chain_id": "f8c74ccb7f9dea6f26a6d7f786809ddd1bce9fada3867f567dd83691b5348534"
                        },
                        "time": 1660745796190
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "DOWN",
                        "service_data": {
                            "last_indexed_block": 0,
                            "total_indexed_blocks": 0,
                            "active_shards": "0.0%"
                        },
                        "time": 1660745796204
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 500)
    check, status := api.Call()

    assert.Equal(t, "Failed to get Elasticsearch and/or nodeos block numbers (es: 0, eos: 263148621)", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2NodeosRPCFailed(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "DOWN",
                        "service_data": {
                            "head_block_num": 0,
                            "head_block_time": "",
                            "time_offset": 0,
                            "last_irreversible_block": 0,
                            "chain_id": ""
                        },
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "DOWN",
                        "service_data": {
                            "last_indexed_block": 263148121,
                            "total_indexed_blocks": 263148121,
                            "active_shards": "100.0%"
                        },
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 500)
    check, status := api.Call()

    assert.Equal(t, "Failed to get Elasticsearch and/or nodeos block numbers (es: 263148121, eos: 0)", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
    assert.Equal(t, expected, check)
}

func TestEosioV2ElasticsNodeosRPCFailed(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v2/health" {
            info := `{
                "version": "1.0",
                "version_hash": "028d5a34463884fcbe2ecfd3c0fcb3b5d4d538f4fd64803c1ef7209c85f2f266",
                "host": "api.test.com:443",
                "health": [
                    {
                        "service": "NodeosRPC",
                        "status": "DOWN",
                        "service_data": {},
                        "time": 1642174781678
                    },
                    {
                        "service": "Elasticsearch",
                        "status": "DOWN",
                        "service_data": {},
                        "time": 1642174781736
                    }
                ]
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV2(srv.URL, "", 500)
    check, status := api.Call()

    assert.Equal(t, "Failed to get Elasticsearch and/or nodeos block numbers (es: 0, eos: 0)", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
    assert.Equal(t, expected, check)
}
