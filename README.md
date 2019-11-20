# protobufc

[![Build Status](https://travis-ci.org/cgrates/protobufc.svg)](https://travis-ci.org/cgrates/protobufc)
[![GoDoc](https://godoc.org/github.com/cgrates/protobufc?status.svg)](https://godoc.org/github.com/cgrates/protobufc)



# Install

Install `protorpc` package:

1. `go get github.com/cgrates/protobufc`
1. `go run hello.go`

Install `protoc-gen-go` plugin:

1. install `protoc` at first: http://github.com/google/protobuf/releases
1. `go get github.com/golang/protobuf/protoc-gen-go`
1. `go get github.com/cgrates/protobufc/protoc-gen-protorpc`
1. `go generate github.com/cgrates/protobufc/examples/service.pb`
1. `go test github.com/cgrates/protobufc/examples/service.pb`


# Examples

First, create [echo.proto](examples/service.pb/echo.proto):

```Proto
syntax = "proto3";

package service;

message EchoRequest {
	string msg = 1;
}

message EchoResponse {
	string msg = 1;
}

service EchoService {
	rpc Echo (EchoRequest) returns (EchoResponse);
	rpc EchoTwice (EchoRequest) returns (EchoResponse);
}
```

Second, generate [echo.pb.go](examples/service.pb/echo.pb.go) and [echo.pb.protorpc.go](examples/service.pb/echo.pb.protorpc.go)
from [echo.proto](examples/service.pb/echo.proto) (we can use `go generate` to invoke this command, see [proto.go](examples/service.pb/proto.go)).

	protoc --go_out=. echo.proto
	protoc --protorpc_out=. echo.proto


Now, we can use the stub code like this:

```Go
package main

import (
	"fmt"
	"log"

	"github.com/cgrates/protobufc"
	service "github.com/cgrates/protobufc/examples/service.pb"
)

type Echo int

func (t *Echo) Echo(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = args.Msg
	return nil
}

func (t *Echo) EchoTwice(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = args.Msg + args.Msg
	return nil
}

func init() {
	go service.ListenAndServeEchoService("tcp", `127.0.0.1:9527`, new(Echo))
}

func main() {
	echoClient, err := service.DialEchoService("tcp", `127.0.0.1:9527`)
	if err != nil {
		log.Fatalf("service.DialEchoService: %v", err)
	}
	defer echoClient.Close()

	args := &service.EchoRequest{Msg: "你好, 世界!"}
	reply, err := echoClient.EchoTwice(args)
	if err != nil {
		log.Fatalf("echoClient.EchoTwice: %v", err)
	}
	fmt.Println(reply.Msg)

	// or use normal client
	client, err := protorpc.Dial("tcp", `127.0.0.1:9527`)
	if err != nil {
		log.Fatalf("protorpc.Dial: %v", err)
	}
	defer client.Close()

	echoClient1 := &service.EchoServiceClient{client}
	echoClient2 := &service.EchoServiceClient{client}
	reply, err = echoClient1.EchoTwice(args)
	reply, err = echoClient2.EchoTwice(args)
	_, _ = reply, err

	// Output:
	// 你好, 世界!你好, 世界!
}
```

[More examples](examples).

# standard net/rpc

First, create [echo.proto](examples/stdrpc.pb/echo.proto):

```Proto
syntax = "proto3";

package service;

message EchoRequest {
	string msg = 1;
}

message EchoResponse {
	string msg = 1;
}

service EchoService {
	rpc Echo (EchoRequest) returns (EchoResponse);
	rpc EchoTwice (EchoRequest) returns (EchoResponse);
}
```

Second, generate [echo.pb.go](examples/stdrpc.pb/echo.pb.go) from [echo.proto](examples/stdrpc.pb/echo.proto) with `protoc-gen-stdrpc` plugin.

	protoc --stdrpc_out=. echo.proto

The stdrpc plugin generated code do not depends **protorpc** package, it use gob as the default rpc encoding.

# BUGS

Report bugs to <chaishushan@gmail.com>.

Thanks!
