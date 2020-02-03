
package eosapi;

import "time";

// get_info format (not all fields).
type Info struct {
	ServerVersion 	string 		`json:"server_version"`
	HeadBlockNum 	int64 		`json:"head_block_num"`
	HeadBlockTime 	time.Time 	`json:"head_block_time"`
}
