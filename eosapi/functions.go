
package eosapi;

import (
	"time"
	"net/url"
	"io/ioutil"
	"github.com/imroc/req"
	"github.com/liamylian/jsontime/v2"
)

var json = v2.ConfigWithCustomTimeFormat

func init() {

	// EOS Api does not specify timezone in timestamps (they are always UTC tho).
	v2.SetDefaultTimeFormat("2006-01-02T15:04:05", time.UTC)
}

func send(method string, api_url string) (*req.Resp, error) {

	u, err := url.Parse(api_url)
    if err != nil {
        return nil, err
    }

	// Go's net.http (that `req` uses) sends the port in the host header.
	// nodeos api does not like that, so we need to provide our
	// own Host header with just the host.
	headers := req.Header{
		"Host": u.Host,
	}

	r := req.New()
	return r.Do(method, api_url, headers)
}

//  GetInfo - Fetches get_info from API
// ---------------------------------------------------------
func GetInfo(url string) (Info, error) {

	var info Info

	r, err := send("GET", url + "/v1/chain/get_info")
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body)

		// Parse json
		err = json.Unmarshal(body, &info)
	}
	return info, err
}

func GetHealth(url string) (Health, error) {

	var health Health;

	r, err := send("GET", url + "/v2/health")
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body)

		// Parse json
		err = json.Unmarshal(body, &health)
	}
	return health, err
}
