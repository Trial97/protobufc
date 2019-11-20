// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	protobufc "github.com/cgrates/protobufc"
	message "github.com/cgrates/protobufc/examples/message.pb"
)

type Echo int

func (t *Echo) Echo(args *message.EchoRequest, reply *message.EchoResponse) error {
	reply.Msg = args.Msg
	return nil
}

var serverStarted chan struct{}

func ListenAndServeEchoService() error {
	listener, err := net.Listen("tcp", ":12347")
	if err != nil {
		return err
	}
	rpcServer := rpc.NewServer()
	if err := rpcServer.Register(new(Echo)); err != nil {
		return err
	}
	serverStarted <- struct{}{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go rpcServer.ServeCodec(protobufc.NewServerCodec(conn))
	}
}

func init() {
	serverStarted = make(chan struct{})
	go func() {
		err := ListenAndServeEchoService()
		log.Fatal(err)
	}()
}

func main() {
	<-serverStarted // wait to start the server
	client, err := protobufc.Dial("tcp", ":12347")
	if err != nil {
		log.Fatalf("protobufc.Dial: %v", err)
	}
	defer client.Close()
	args := &message.EchoRequest{Msg: "Hello world!"}
	reply := new(message.EchoResponse)
	if err := client.Call("Echo.Echo", args, reply); err != nil {
		log.Fatalf("protobufc.Call: %v", err)
	}
	fmt.Println(reply.Msg)
}
