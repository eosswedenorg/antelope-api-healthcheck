
package api

import (
    "time"
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
    "github.com/eosswedenorg-go/haproxy/agentcheck"
)

func TestV1LogInfo(t *testing.T) {

    api := NewEosioV1("https://api.v1.example.com", "host.example.com", 120)

    expected := LogParams{"type","eosio-v1","url","https://api.v1.example.com","host","host.example.com","block_time",float64(120)}

    assert.Equal(t, expected, api.LogInfo())
}

func TestV1SetTime(t *testing.T) {

    expected := time.Date(2022, 2, 24, 13, 38, 0, 0, time.UTC)

    api := NewEosioV1("", "", 60)
    // Assert that time is NOW (+-10 seconds)
    assert.InDelta(t, api.GetTime().Unix(), time.Now().In(time.UTC).Unix(), float64(10))

    api.SetTime(expected)
    assert.Equal(t, expected, api.GetTime())
}

func TestV1JsonFailure(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.Write([]byte(`!//{invalid-json}!##`))
    }))

    api := NewEosioV1(srv.URL, "", 120)
    check, _ := api.Call()

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Failed, "")
    assert.Equal(t, expected, check)
}

func TestV1HTTP500Down(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(500)
        res.Write([]byte(`{}`))
    }))

    api := NewEosioV1(srv.URL, "", 120)
    check, status := api.Call()

    assert.Equal(t, "Taking offline because 500 was received from backend", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}

func TestV1LaggingUp(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v1/chain/get_info" {
            info := `{
                "server_version": "8f613ec9",
                "head_block_num": 7272812,
                "head_block_time": "2022-02-24T13:37:00"
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV1(srv.URL, "", 60)
    api.SetTime(time.Date(2022, 2, 24, 13, 38, 0, 0, time.UTC))
    check, status := api.Call()

    assert.Equal(t, "OK", status)

    expected := agentcheck.NewStatusResponse(agentcheck.Up)
    assert.Equal(t, expected, check)
}

func TestV1LaggingDown(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v1/chain/get_info" {
            info := `{
                "server_version": "9a607cce",
                "head_block_num": 87263,
                "head_block_time": "2018-01-01T13:37:01"
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV1(srv.URL, "", 60)
    api.SetTime(time.Date(2018, time.January, 1, 13, 38, 2, 0, time.UTC))
    check, status := api.Call()

    assert.Equal(t, "Taking offline because head block is lagging 61 seconds", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}

func TestV1TimeInFutureUP(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v1/chain/get_info" {
            info := `{
                "server_version": "d1bec8d3",
                "head_block_num": 548847,
                "head_block_time": "2020-09-22T09:32:00"
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV1(srv.URL, "", 120)
    api.SetTime(time.Date(2020, 9, 22, 9, 30, 0, 0, time.UTC))
    check, status := api.Call()

    assert.Equal(t, "OK", status)

    expected := agentcheck.NewStatusResponse(agentcheck.Up)
    assert.Equal(t, expected, check)
}


func TestV1TimeInFutureDown(t *testing.T) {

    var srv = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        if req.URL.String() == "/v1/chain/get_info" {
            info := `{
                "server_version": "c879d231",
                "head_block_num": 2637621,
                "head_block_time": "2019-04-14T12:02:01"
            }`

            res.Write([]byte(info))
        }
    }))

    api := NewEosioV1(srv.URL, "", 120)
    api.SetTime(time.Date(2019, time.April, 14, 12, 0, 0, 0, time.UTC))
    check, status := api.Call()

    assert.Equal(t, "Taking offline because head block is -121 seconds into the future", status)

    expected := agentcheck.NewStatusMessageResponse(agentcheck.Down, "")
    assert.Equal(t, expected, check)
}
