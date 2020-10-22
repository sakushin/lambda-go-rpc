package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
)

func main() {
	var (
		host string
		port string
	)
	flag.StringVar(&host, "host", "127.0.0.1", "target host")
	flag.StringVar(&port, "port", "", "target port")
	flag.Parse()

	if len(port) == 0 {
		port = os.Getenv("_LAMBDA_SERVER_PORT")
		if len(port) == 0 {
			errprintln("-port required")
			flag.Usage()
			os.Exit(1)
		}
	}

	payload, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	addr := host + ":" + port
	conn, err := rpc.Dial("tcp", addr)
	if err != nil {
		errprintln(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	req := request{Payload: payload}
	var res response
	if err := conn.Call("Function.Invoke", &req, &res); err != nil {
		errprintln(fmt.Sprintf("invokation error: %s", err.Error()))
		os.Exit(1)
	}

	if res.Error != nil {
		errprintln(res.Error.Message)
		os.Exit(1)
	}

	fmt.Print(string(res.Payload))
}

func errprintln(message string) {
	fmt.Fprintln(os.Stderr, message)
}

type request struct {
	Payload []byte
}

type response struct {
	Payload []byte
	Error   *responseError
}

type responseError struct {
	Message string `json:"errorMessage"`
}
