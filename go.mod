module github.com/eosswedenorg/eosio-api-healthcheck

go 1.14

require (
	github.com/eosswedenorg-go/haproxy v0.0.0-20220101140534-fccfdd93a8cd
	github.com/eosswedenorg-go/pid v1.0.0
	github.com/eosswedenorg-go/tcp_server v0.1.0
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/inconshreveable/log15 v0.0.0-20201112154412-8562bdadbbac
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/pborman/getopt/v2 v2.1.0
	internal/eosapi v1.0.0
	internal/utils v1.0.0
)

replace internal/eosapi => ./src/eosapi

replace internal/utils => ./src/utils
