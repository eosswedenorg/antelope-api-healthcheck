
package eosapi;

import (
    "fmt"
	"time"
	"io/ioutil"
	"github.com/imroc/req"
	"github.com/liamylian/jsontime/v2"
)

var json = v2.ConfigWithCustomTimeFormat;

func init() {

	// EOS Api does not specify timezone in timestamps (they are always UTC tho).
	v2.SetDefaultTimeFormat("2006-01-02T15:04:05", time.UTC);
}

func send(method string, host string, port int, uri string) (*req.Resp, error) {

	// Go's net.http (that `req` uses) sends the port in the host header.
	// nodeos api does not like that, so we need to provide our
	// own Host header with just the host.
	headers := req.Header{
		"Host": host,
	}

	r := req.New()
	return r.Do(method, fmt.Sprintf("http://%s:%d%s", host, port, uri), headers);
}

//  GetInfo - Fetches get_info from API
// ---------------------------------------------------------
func GetInfo(host string, port int) (Info, error) {

	var info Info;

	r, err := send("GET", host, port, "/v1/chain/get_info");
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body);

		// Parse json
		err = json.Unmarshal(body, &info);
	}
	return info, err;
}

func GetHealth(host string, port int) (Health, error) {

	var health Health;

	r, err := send("GET", host, port, "/v2/health");
	if err == nil {
		resp := r.Response()
		body, _ := ioutil.ReadAll(resp.Body);

		// Parse json
		err = json.Unmarshal(body, &health);
	}
	return health, err;
}
