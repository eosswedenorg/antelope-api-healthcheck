
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

//  GetInfo - Fetches get_info from API
// ---------------------------------------------------------
func GetInfo(host string, port int) (Info, error) {

	var info Info;

	// Format url.
	url := fmt.Sprintf("http://%s:%d/v1/chain/get_info", host, port);

	// Go's net.http (that `req` uses) sends the port in the host header.
	// nodeos api does not like that, so we need to provide our
	// own Host header with just the host.
	headers := req.Header{
		"Host": host,
	}

	// Send HTTP Get request.
	r, err := req.Get(url, headers)
	resp := r.Response()
	body, _ := ioutil.ReadAll(resp.Body);

	// Parse json
	err = json.Unmarshal(body, &info);
	return info, err;
}
