// Code generated by protoc-gen-stdrpc. DO NOT EDIT.
//
// plugin: https://github.com/chai2010/protorpc/tree/master/protoc-gen-stdrpc
// plugin: https://github.com/chai2010/protorpc/tree/master/protoc-plugin-common
//
// source: echo.proto

package service

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	_ = fmt.Sprint
	_ = io.Reader(nil)
	_ = log.Print
	_ = net.Addr(nil)
	_ = rpc.Call{}
	_ = time.Second

	_ = proto.String
)

type EchoService interface {
	Echo(in *EchoRequest, out *EchoResponse) error
	EchoTwice(in *EchoRequest, out *EchoResponse) error
}

// AcceptEchoServiceClient accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func AcceptEchoServiceClient(lis net.Listener, x EchoService) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeConn(conn)
	}
}

// RegisterEchoService publish the given EchoService implementation on the server.
func RegisterEchoService(srv *rpc.Server, x EchoService) error {
	if err := srv.RegisterName("EchoService", x); err != nil {
		return err
	}
	return nil
}

// NewEchoServiceServer returns a new EchoService Server.
func NewEchoServiceServer(x EchoService) *rpc.Server {
	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		log.Fatal(err)
	}
	return srv
}

// ListenAndServeEchoService listen announces on the local network address laddr
// and serves the given EchoService implementation.
func ListenAndServeEchoService(network, addr string, x EchoService) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("EchoService", x); err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeConn(conn)
	}
}

type EchoServiceClient struct {
	*rpc.Client
}

// NewEchoServiceClient returns a EchoService stub to handle
// requests to the set of EchoService at the other end of the connection.
func NewEchoServiceClient(conn io.ReadWriteCloser) *EchoServiceClient {
	c := rpc.NewClient(conn)
	return &EchoServiceClient{c}
}

func (c *EchoServiceClient) Echo(in *EchoRequest) (out *EchoResponse, err error) {
	if in == nil {
		in = new(EchoRequest)
	}

	type Validator interface {
		Validate() error
	}
	if x, ok := proto.Message(in).(Validator); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	out = new(EchoResponse)
	if err = c.Call("EchoService.Echo", in, out); err != nil {
		return nil, err
	}

	if x, ok := proto.Message(out).(Validator); ok {
		if err := x.Validate(); err != nil {
			return out, err
		}
	}

	return out, nil
}

func (c *EchoServiceClient) AsyncEcho(in *EchoRequest, out *EchoResponse, done chan *rpc.Call) *rpc.Call {
	if in == nil {
		in = new(EchoRequest)
	}
	return c.Go(
		"EchoService.Echo",
		in, out,
		done,
	)
}

func (c *EchoServiceClient) EchoTwice(in *EchoRequest) (out *EchoResponse, err error) {
	if in == nil {
		in = new(EchoRequest)
	}

	type Validator interface {
		Validate() error
	}
	if x, ok := proto.Message(in).(Validator); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	out = new(EchoResponse)
	if err = c.Call("EchoService.EchoTwice", in, out); err != nil {
		return nil, err
	}

	if x, ok := proto.Message(out).(Validator); ok {
		if err := x.Validate(); err != nil {
			return out, err
		}
	}

	return out, nil
}

func (c *EchoServiceClient) AsyncEchoTwice(in *EchoRequest, out *EchoResponse, done chan *rpc.Call) *rpc.Call {
	if in == nil {
		in = new(EchoRequest)
	}
	return c.Go(
		"EchoService.EchoTwice",
		in, out,
		done,
	)
}

// DialEchoService connects to an EchoService at the specified network address.
func DialEchoService(network, addr string) (*EchoServiceClient, error) {
	c, err := rpc.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &EchoServiceClient{c}, nil
}

// DialEchoServiceTimeout connects to an EchoService at the specified network address.
func DialEchoServiceTimeout(network, addr string, timeout time.Duration) (*EchoServiceClient, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return &EchoServiceClient{rpc.NewClient(conn)}, nil
}