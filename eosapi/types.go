
package eosapi;

import "time"

// get_info format (not all fields).
type Info struct {
	ServerVersion 	string 		`json:"server_version"`
	HeadBlockNum 	int64 		`json:"head_block_num"`
	HeadBlockTime 	time.Time 	`json:"head_block_time"`
}

// Service struct from /v2/health
type Service struct {
	Name	string 				   `json:"service"`
	Status  string  			   `json:"status"`
	Data    map[string]interface{} `json:"service_data"`
	Time    int64 				   `json:"time"` // unix timestamp.
}

// /v2/health format (not all fields).
type Health struct {
	VersionHash 	string 		`json:"version_hash"`
	Health          []Service 	`json:"health"`
}
