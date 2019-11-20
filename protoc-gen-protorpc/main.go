// Copyright 2018 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"text/template"

	plugin "github.com/cgrates/protobufc/protoc-gen-plugin"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/generator"
)

func main() {
	plugin.Main()
}

func init() {
	plugin.RegisterCodeGenerator(new(protorpcPlugin))
}

type protorpcPlugin struct{}

func (p *protorpcPlugin) Name() string        { return "protorpc-go" }
func (p *protorpcPlugin) FileNameExt() string { return ".pb.protorpc.go" }

func (p *protorpcPlugin) HeaderCode(g *generator.Generator, file *generator.FileDescriptor) string {
	const tmpl = `
{{- $G := .G -}}
{{- $File := .File -}}

// Code generated by protoc-gen-protorpc. DO NOT EDIT.
//
// plugin: https://github.com/cgrates/protobufc/tree/master/protoc-gen-plugin
// plugin: https://github.com/cgrates/protobufc/tree/master/protoc-gen-protorpc
//
// source: {{$File.GetName}}

package {{$File.PackageName}}

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/cgrates/protobufc"
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
	_ = protorpc.Dial
)
`
	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmpl))
	err := t.Execute(&buf,
		struct {
			G    *generator.Generator
			File *generator.FileDescriptor
		}{
			G:    g,
			File: file,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func (p *protorpcPlugin) ServiceCode(g *generator.Generator, file *generator.FileDescriptor, svc *descriptor.ServiceDescriptorProto) string {
	var code string
	code += p.genServiceInterface(g, file, svc)
	code += p.genServiceServer(g, file, svc)
	code += p.genServiceClient(g, file, svc)
	return code
}

func (p *protorpcPlugin) MessageCode(g *generator.Generator, file *generator.FileDescriptor, msg *descriptor.DescriptorProto) string {
	return ""
}

func (p *protorpcPlugin) genServiceInterface(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const serviceInterfaceTmpl = `
type {{.ServiceName}} interface {
	{{.CallMethodList}}
}
`
	const callMethodTmpl = `
{{.MethodName}}(in *{{.ArgsType}}, out *{{.ReplyType}}) error`

	// gen call method list
	var callMethodList string
	for _, m := range svc.Method {
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(callMethodTmpl))
		t.Execute(out, &struct{ ServiceName, MethodName, ArgsType, ReplyType string }{
			ServiceName: generator.CamelCase(svc.GetName()),
			MethodName:  generator.CamelCase(m.GetName()),
			ArgsType:    g.TypeName(g.ObjectNamed(m.GetInputType())),
			ReplyType:   g.TypeName(g.ObjectNamed(m.GetOutputType())),
		})
		callMethodList += out.String()

		g.RecordTypeUse(m.GetInputType())
		g.RecordTypeUse(m.GetOutputType())
	}

	// gen all interface code
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(serviceInterfaceTmpl))
		t.Execute(out, &struct{ ServiceName, CallMethodList string }{
			ServiceName:    generator.CamelCase(svc.GetName()),
			CallMethodList: callMethodList,
		})

		return out.String()
	}
}

func (p *protorpcPlugin) genServiceServer(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const serviceHelperFunTmpl = `
// Accept{{.ServiceName}}Client accepts connections on the listener and serves requests
// for each incoming connection.  Accept blocks; the caller typically
// invokes it in a go statement.
func Accept{{.ServiceName}}Client(lis net.Listener, x {{.ServiceName}}) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// Register{{.ServiceName}} publish the given {{.ServiceName}} implementation on the server.
func Register{{.ServiceName}}(srv *rpc.Server, x {{.ServiceName}}) error {
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		return err
	}
	return nil
}

// New{{.ServiceName}}Server returns a new {{.ServiceName}} Server.
func New{{.ServiceName}}Server(x {{.ServiceName}}) *rpc.Server {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}
	return srv
}

// ListenAndServe{{.ServiceName}} listen announces on the local network address laddr
// and serves the given {{.ServiceName}} implementation.
func ListenAndServe{{.ServiceName}}(network, addr string, x {{.ServiceName}}) error {
	lis, err := net.Listen(network, addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		return err
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalf("lis.Accept(): %v\n", err)
		}
		go srv.ServeCodec(protorpc.NewServerCodec(conn))
	}
}

// Serve{{.ServiceName}} serves the given {{.ServiceName}} implementation.
func Serve{{.ServiceName}}(conn io.ReadWriteCloser, x {{.ServiceName}}) {
	srv := rpc.NewServer()
	if err := srv.RegisterName("{{.ServiceRegisterName}}", x); err != nil {
		log.Fatal(err)
	}
	srv.ServeCodec(protorpc.NewServerCodec(conn))
}
`
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(serviceHelperFunTmpl))
		t.Execute(out, &struct{ PackageName, ServiceName, ServiceRegisterName string }{
			PackageName: file.GetPackage(),
			ServiceName: generator.CamelCase(svc.GetName()),
			ServiceRegisterName: p.makeServiceRegisterName(
				file, file.GetPackage(), generator.CamelCase(svc.GetName()),
			),
		})

		return out.String()
	}
}

func (p *protorpcPlugin) genServiceClient(
	g *generator.Generator,
	file *generator.FileDescriptor,
	svc *descriptor.ServiceDescriptorProto,
) string {
	const clientHelperFuncTmpl = `
type {{.ServiceName}}Client struct {
	*rpc.Client
}

// New{{.ServiceName}}Client returns a {{.ServiceName}} stub to handle
// requests to the set of {{.ServiceName}} at the other end of the connection.
func New{{.ServiceName}}Client(conn io.ReadWriteCloser) (*{{.ServiceName}}Client) {
	c := rpc.NewClientWithCodec(protorpc.NewClientCodec(conn))
	return &{{.ServiceName}}Client{c}
}

{{.MethodList}}

// Dial{{.ServiceName}} connects to an {{.ServiceName}} at the specified network address.
func Dial{{.ServiceName}}(network, addr string) (*{{.ServiceName}}Client, error) {
	c, err := protorpc.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return &{{.ServiceName}}Client{c}, nil
}

// Dial{{.ServiceName}}Timeout connects to an {{.ServiceName}} at the specified network address.
func Dial{{.ServiceName}}Timeout(network, addr string, timeout time.Duration) (*{{.ServiceName}}Client, error) {
	c, err := protorpc.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return &{{.ServiceName}}Client{c}, nil
}
`
	const clientMethodTmpl = `
func (c *{{.ServiceName}}Client) {{.MethodName}}(in *{{.ArgsType}}) (out *{{.ReplyType}}, err error) {
	if in == nil {
		in = new({{.ArgsType}})
	}

	type Validator interface {
		Validate() error
	}
	if x, ok := proto.Message(in).(Validator); ok {
		if err := x.Validate(); err != nil {
			return nil, err
		}
	}

	out = new({{.ReplyType}})
	if err = c.Call("{{.ServiceRegisterName}}.{{.MethodName}}", in, out); err != nil {
		return nil, err
	}

	if x, ok := proto.Message(out).(Validator); ok {
		if err := x.Validate(); err != nil {
			return out, err
		}
	}

	return out, nil
}

func (c *{{.ServiceName}}Client) Async{{.MethodName}}(in *{{.ArgsType}}, out *{{.ReplyType}}, done chan *rpc.Call) *rpc.Call {
	if in == nil {
		in = new({{.ArgsType}})
	}
	return c.Go(
		"{{.ServiceRegisterName}}.{{.MethodName}}",
		in, out,
		done,
	)
}
`

	// gen client method list
	var methodList string
	for _, m := range svc.Method {
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(clientMethodTmpl))
		t.Execute(out, &struct{ ServiceName, ServiceRegisterName, MethodName, ArgsType, ReplyType string }{
			ServiceName: generator.CamelCase(svc.GetName()),
			ServiceRegisterName: p.makeServiceRegisterName(
				file, file.GetPackage(), generator.CamelCase(svc.GetName()),
			),
			MethodName: generator.CamelCase(m.GetName()),
			ArgsType:   g.TypeName(g.ObjectNamed(m.GetInputType())),
			ReplyType:  g.TypeName(g.ObjectNamed(m.GetOutputType())),
		})
		methodList += out.String()
	}

	// gen all client code
	{
		out := bytes.NewBuffer([]byte{})
		t := template.Must(template.New("").Parse(clientHelperFuncTmpl))
		t.Execute(out, &struct{ PackageName, ServiceName, MethodList string }{
			PackageName: file.GetPackage(),
			ServiceName: generator.CamelCase(svc.GetName()),
			MethodList:  methodList,
		})

		return out.String()
	}
}

func (p *protorpcPlugin) makeServiceRegisterName(
	file *generator.FileDescriptor,
	packageName, serviceName string,
) string {
	// return packageName + "." + serviceName
	return serviceName
}
