// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protobufc

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"sync"

	wire "github.com/cgrates/protobufc/wire.pb"
	"github.com/gogo/protobuf/proto"
)

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	// temporary work space
	reqHeader wire.RequestHeader

	// Package rpc expects uint64 request IDs.
	// We assign uint64 sequence numbers to incoming requests
	// but save the original request ID in the pending map.
	// When rpc responds, we use the sequence number in
	// the response to find the original request ID.
	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]uint64
}

// NewServerCodec returns a serverCodec that communicates with the ClientCodec
// on the other end of the given conn.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{
		r:       conn,
		w:       conn,
		c:       conn,
		pending: make(map[uint64]uint64),
	}
}

func (s *serverCodec) ReadRequestHeader(r *rpc.Request) (err error) {
	header := wire.RequestHeader{}
	if err = readRequestHeader(s.r, &header); err != nil {
		return
	}

	s.mutex.Lock()
	s.seq++
	s.pending[s.seq] = header.Id
	s.mutex.Unlock()
	r.ServiceMethod = header.Method
	r.Seq = s.seq

	s.reqHeader = header
	return
}

func (s *serverCodec) ReadRequestBody(x interface{}) (err error) {
	if x == nil {
		return
	}
	request, ok := x.(proto.Message)
	if !ok {
		return fmt.Errorf(
			"protobufc.ServerCodec.ReadRequestBody: %T does not implement proto.Message",
			x,
		)
	}

	if err = readRequestBody(s.r, &s.reqHeader, request); err != nil {
		return
	}

	s.reqHeader = wire.RequestHeader{}
	return
}

func (s *serverCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	var response proto.Message
	if x != nil {
		var ok bool
		if response, ok = x.(proto.Message); !ok {
			if _, ok = x.(struct{}); !ok {
				s.mutex.Lock()
				delete(s.pending, r.Seq)
				s.mutex.Unlock()
				return fmt.Errorf(
					"protobufc.ServerCodec.WriteResponse: %T does not implement proto.Message",
					x,
				)
			}
		}
	}

	s.mutex.Lock()
	id, ok := s.pending[r.Seq]
	if !ok {
		s.mutex.Unlock()
		return errors.New("protobufc: invalid sequence number in response")
	}
	delete(s.pending, r.Seq)
	s.mutex.Unlock()

	return writeResponse(s.w, id, r.Error, response)
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}

// ServeConn runs the Protobuf-RPC server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
func ServeConn(conn io.ReadWriteCloser) {
	rpc.ServeCodec(NewServerCodec(conn))
}
