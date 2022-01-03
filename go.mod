module github.com/eosswedenorg/eosio-api-healthcheck

go 1.14

require (
	internal/eosapi v1.0.0
	internal/log v1.0.0
	internal/utils v1.0.0

	github.com/eosswedenorg-go/pid v1.0.0
	github.com/eosswedenorg-go/haproxy v0.0.0-20220101140534-fccfdd93a8cd
	github.com/firstrow/tcp_server v0.1.0
	github.com/pborman/getopt/v2 v2.1.0
)

replace internal/eosapi => ./src/eosapi

replace internal/log => ./src/log

replace internal/utils => ./src/utils
