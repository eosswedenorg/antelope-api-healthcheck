
package eosapi;

import (
	"strings"
	"time"
	"net/url"
	"io/ioutil"
	"github.com/imroc/req"
	jsontime "github.com/liamylian/jsontime/v2/v2"
)

type ReqParams struct {
	Url string
	Host string
}

var json = jsontime.ConfigWithCustomTimeFormat

func init() {

	// EOS Api does not specify timezone in timestamps (they are always UTC tho).
	jsontime.SetDefaultTimeFormat("2006-01-02T15:04:05", time.UTC)
}

func send(p ReqParams, method string, path string) (*req.Resp, error) {

	host := p.Host
	if len(host) < 1 {
		u, err := url.Parse(p.Url)
		if err != nil {
	        return nil, err
	    }
		host = strings.Split(u.Host, ":")[0]
	}

	// Go's net.http (that `req` uses) sends the port in the host header.
	// nodeos api does not like that, so we need to provide our
	// own Host header with just the host.
	headers := req.Header{
		"Host": host,
	}

	r := req.New()
	return r.Do(method, p.Url + path, headers)
}

//  GetInfo - Fetches get_info from API
// ---------------------------------------------------------
func GetInfo(params ReqParams) (Info, error) {

	var info Info

	r, err := send(params, "GET", "/v1/chain/get_info")
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body)

		// Set HTTPStatusCode
		info.HTTPStatusCode = resp.StatusCode

		// Parse json
		err = json.Unmarshal(body, &info)
	}
	return info, err
}

func GetHealth(params ReqParams) (Health, error) {

	var health Health;

	r, err := send(params, "GET", "/v2/health")
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body)

        // Set HTTPStatusCode
        health.HTTPStatusCode = resp.StatusCode

		// Parse json
		err = json.Unmarshal(body, &health)
	}
	return health, err
}
