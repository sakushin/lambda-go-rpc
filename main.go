package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"time"
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

	deadline := time.Now().Add(15 * time.Minute) // fixed
	req := InvokeRequest{
		Payload: payload,
		Deadline: InvokeRequest_Timestamp{
			Seconds: deadline.Unix(),
			Nanos:   int64(deadline.Nanosecond()),
		},
	}

	var res InvokeResponse
	if err := conn.Call("Function.Invoke", &req, &res); err != nil {
		errprintln(fmt.Sprintf("invocation error: %s", err.Error()))
		os.Exit(1)
	}

	if res.Error != nil {
		errprintln(res.Error.Error())
		os.Exit(1)
	}

	fmt.Print(string(res.Payload))
}

func errprintln(message string) {
	fmt.Fprintln(os.Stderr, message)
}

// see: https://github.com/aws/aws-lambda-go/blob/main/lambda/messages/messages.go

type InvokeRequest_Timestamp struct {
	Seconds int64
	Nanos   int64
}

type InvokeRequest struct {
	Payload               []byte
	RequestId             string
	XAmznTraceId          string
	Deadline              InvokeRequest_Timestamp
	InvokedFunctionArn    string
	CognitoIdentityId     string
	CognitoIdentityPoolId string
	ClientContext         []byte
}

type InvokeResponse struct {
	Payload []byte
	Error   *InvokeResponse_Error
}

type InvokeResponse_Error struct {
	Message    string                             `json:"errorMessage"`
	Type       string                             `json:"errorType"`
	StackTrace []*InvokeResponse_Error_StackFrame `json:"stackTrace,omitempty"`
	ShouldExit bool                               `json:"-"`
}

func (e InvokeResponse_Error) Error() string {
	return fmt.Sprintf("%#v", e)
}

type InvokeResponse_Error_StackFrame struct {
	Path  string `json:"path"`
	Line  int32  `json:"line"`
	Label string `json:"label"`
}
