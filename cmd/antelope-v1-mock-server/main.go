package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/eosswedenorg-go/leapapi"
)

var (
	listen_host = flag.String("h,host", "localhost", "Host to listen on.")
	listen_port = flag.Int("p", 3333, "Port to listen to.")
)

func getInfo(w http.ResponseWriter, r *http.Request) {
	current_time := time.Now()

	info := leapapi.Info{
		ServerVersion:             "c83ea9c2",
		ServerVersionString:       "0.0.0-debug",
		ServerFullVersionString:   "0.0.0-debug-c83ea9c21f60670a00627319ebbd233e6bb4f84904dbcfc894242ba38b2761d4",
		HeadBlockNum:              1000,
		HeadBlockID:               "168d2cf232ca78e94d57a86301e35f110b6016358e05d49ab822df0a8aa988ea",
		HeadBlockTime:             current_time.UTC(),
		ChainID:                   "1045fa26e1c5be590ae6114e73331152671f13c87eee60a2171387dcbc44da88",
		HeadBlockProducer:         "debugproducer",
		LastIrreversableBlockNum:  900,
		LastIrreversableBlockID:   "5149254b9b6fd61a02403ebe3b45ade57642ed473295f33e2184e56966370a1f",
		LastIrreversableBlockTime: current_time.Add(time.Second * -5).UTC(),
		VirtualBlockCPULimit:      4000,
		VirtualBlockNETLimit:      5000,
		BlockCPULimit:             8000,
		BlockNETLimit:             2000,
		TotalCPUWeight:            60488453825414473,
		TotalNETWeight:            101764028077814346,
		ForkDBHeadBlockID:         "7544799d7c2f511368cb94adc65223e1e2cc4cf9639ba07eef2421486a8dbfe5",
		ForkDBHeadBlockNum:        100,
	}

	payload, err := leapapi.Json().Marshal(&info)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	flag.Parse()

	http.HandleFunc("/v1/chain/get_info", getInfo)

	addr := fmt.Sprintf("%s:%d", *listen_host, *listen_port)

	fmt.Println("Listening on:", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
