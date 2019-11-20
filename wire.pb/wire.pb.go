// Code generated by protoc-gen-go. DO NOT EDIT.
// source: wire.proto

/*
Package protobufc_wire is a generated protocol buffer package.


	protobufc wire format wrapper

	0. Frame Format
	len : uvarint64
	data: byte[len]

	1. Client Send Request
	Send RequestHeader: sendFrame(zsock, hdr, len(hdr))
	Send Request: sendFrame(zsock, body, hdr.snappy_compressed_request_len)

	2. Server Recv Request
	Recv RequestHeader: recvFrame(zsock, hdr, max_hdr_len, 0)
	Recv Request: recvFrame(zsock, body, hdr.snappy_compressed_request_len, 0)

	3. Server Send Response
	Send ResponseHeader: sendFrame(zsock, hdr, len(hdr))
	Send Response: sendFrame(zsock, body, hdr.snappy_compressed_response_len)

	4. Client Recv Response
	Recv ResponseHeader: recvFrame(zsock, hdr, max_hdr_len, 0)
	Recv Response: recvFrame(zsock, body, hdr.snappy_compressed_response_len, 0)


It is generated from these files:
	wire.proto

It has these top-level messages:
	RequestHeader
	ResponseHeader
*/
package protobufc_wire

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Const int32

const (
	Const_ZERO                   Const = 0
	Const_MAX_REQUEST_HEADER_LEN Const = 1024
)

var Const_name = map[int32]string{
	0:    "ZERO",
	1024: "MAX_REQUEST_HEADER_LEN",
}
var Const_value = map[string]int32{
	"ZERO":                   0,
	"MAX_REQUEST_HEADER_LEN": 1024,
}

func (x Const) String() string {
	return proto.EnumName(Const_name, int32(x))
}
func (Const) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type RequestHeader struct {
	Id                         uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Method                     string `protobuf:"bytes,2,opt,name=method" json:"method,omitempty"`
	RawRequestLen              uint32 `protobuf:"varint,3,opt,name=raw_request_len,json=rawRequestLen" json:"raw_request_len,omitempty"`
	SnappyCompressedRequestLen uint32 `protobuf:"varint,4,opt,name=snappy_compressed_request_len,json=snappyCompressedRequestLen" json:"snappy_compressed_request_len,omitempty"`
	Checksum                   uint32 `protobuf:"varint,5,opt,name=checksum" json:"checksum,omitempty"`
}

func (m *RequestHeader) Reset()                    { *m = RequestHeader{} }
func (m *RequestHeader) String() string            { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()               {}
func (*RequestHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RequestHeader) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RequestHeader) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *RequestHeader) GetRawRequestLen() uint32 {
	if m != nil {
		return m.RawRequestLen
	}
	return 0
}

func (m *RequestHeader) GetSnappyCompressedRequestLen() uint32 {
	if m != nil {
		return m.SnappyCompressedRequestLen
	}
	return 0
}

func (m *RequestHeader) GetChecksum() uint32 {
	if m != nil {
		return m.Checksum
	}
	return 0
}

type ResponseHeader struct {
	Id                          uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Error                       string `protobuf:"bytes,2,opt,name=error" json:"error,omitempty"`
	RawResponseLen              uint32 `protobuf:"varint,3,opt,name=raw_response_len,json=rawResponseLen" json:"raw_response_len,omitempty"`
	SnappyCompressedResponseLen uint32 `protobuf:"varint,4,opt,name=snappy_compressed_response_len,json=snappyCompressedResponseLen" json:"snappy_compressed_response_len,omitempty"`
	Checksum                    uint32 `protobuf:"varint,5,opt,name=checksum" json:"checksum,omitempty"`
}

func (m *ResponseHeader) Reset()                    { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string            { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()               {}
func (*ResponseHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ResponseHeader) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ResponseHeader) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *ResponseHeader) GetRawResponseLen() uint32 {
	if m != nil {
		return m.RawResponseLen
	}
	return 0
}

func (m *ResponseHeader) GetSnappyCompressedResponseLen() uint32 {
	if m != nil {
		return m.SnappyCompressedResponseLen
	}
	return 0
}

func (m *ResponseHeader) GetChecksum() uint32 {
	if m != nil {
		return m.Checksum
	}
	return 0
}

func init() {
	proto.RegisterType((*RequestHeader)(nil), "protorpc.wire.RequestHeader")
	proto.RegisterType((*ResponseHeader)(nil), "protorpc.wire.ResponseHeader")
	proto.RegisterEnum("protorpc.wire.Const", Const_name, Const_value)
}

func init() { proto.RegisterFile("wire.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x5f, 0x4b, 0xf3, 0x30,
	0x14, 0xc6, 0xdf, 0xec, 0xdd, 0xc6, 0x3c, 0xd0, 0x5a, 0x82, 0x8c, 0xb2, 0xa1, 0x94, 0x5d, 0x48,
	0xf1, 0xa2, 0x37, 0x7e, 0x82, 0x52, 0x03, 0xbb, 0x98, 0x8a, 0x51, 0x41, 0xbc, 0x09, 0xb5, 0x3d,
	0xb0, 0xa2, 0x6d, 0x62, 0xd2, 0x51, 0xbc, 0xf3, 0x93, 0x09, 0x7e, 0x33, 0x59, 0x53, 0x66, 0xc5,
	0x3f, 0x57, 0xe1, 0x1c, 0x9e, 0xdf, 0x93, 0xfc, 0x08, 0x40, 0x53, 0x68, 0x8c, 0x94, 0x96, 0xb5,
	0xa4, 0x4e, 0x7b, 0x68, 0x95, 0x45, 0xdb, 0xe5, 0xe2, 0x8d, 0x80, 0xc3, 0xf1, 0x79, 0x83, 0xa6,
	0x5e, 0x62, 0x9a, 0xa3, 0xa6, 0x2e, 0x0c, 0x8a, 0xdc, 0x27, 0x01, 0x09, 0x87, 0x7c, 0x50, 0xe4,
	0x74, 0x0a, 0xe3, 0x12, 0xeb, 0xb5, 0xcc, 0xfd, 0x41, 0x40, 0xc2, 0x3d, 0xde, 0x4d, 0xf4, 0x18,
	0xf6, 0x75, 0xda, 0x08, 0x6d, 0x61, 0xf1, 0x84, 0x95, 0xff, 0x3f, 0x20, 0xa1, 0xc3, 0x1d, 0x9d,
	0x36, 0x5d, 0xe5, 0x0a, 0x2b, 0x1a, 0xc3, 0xa1, 0xa9, 0x52, 0xa5, 0x5e, 0x44, 0x26, 0x4b, 0xa5,
	0xd1, 0x18, 0xcc, 0xbf, 0x50, 0xc3, 0x96, 0x9a, 0xd9, 0x50, 0xb2, 0xcb, 0xf4, 0x2a, 0x66, 0x30,
	0xc9, 0xd6, 0x98, 0x3d, 0x9a, 0x4d, 0xe9, 0x8f, 0xda, 0xf4, 0x6e, 0x5e, 0xbc, 0x13, 0x70, 0x39,
	0x1a, 0x25, 0x2b, 0x83, 0xbf, 0x18, 0x1c, 0xc0, 0x08, 0xb5, 0x96, 0xba, 0x13, 0xb0, 0x03, 0x0d,
	0xc1, 0xb3, 0xef, 0xb7, 0x6c, 0x4f, 0xc0, 0x6d, 0x05, 0xec, 0x7a, 0x7b, 0x7d, 0x02, 0x47, 0x3f,
	0x19, 0xf4, 0x38, 0xab, 0x30, 0xff, 0xae, 0xf0, 0x59, 0xf2, 0x87, 0xc3, 0x49, 0x04, 0xa3, 0x44,
	0x56, 0xa6, 0xa6, 0x13, 0x18, 0xde, 0x33, 0x7e, 0xe9, 0xfd, 0xa3, 0x73, 0x98, 0x9e, 0xc7, 0x77,
	0x82, 0xb3, 0xab, 0x5b, 0x76, 0x7d, 0x23, 0x96, 0x2c, 0x3e, 0x63, 0x5c, 0xac, 0xd8, 0x85, 0xf7,
	0x3a, 0x79, 0x18, 0xb7, 0x7f, 0x78, 0xfa, 0x11, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x10, 0xb3, 0xa4,
	0xd8, 0x01, 0x00, 0x00,
}
