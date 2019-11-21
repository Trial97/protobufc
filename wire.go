// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protobufc

import (
	"fmt"
	"io"

	wire "github.com/cgrates/protobufc/wire.pb"
	"github.com/gogo/protobuf/proto"
)

func writeRequest(w io.Writer, id uint64, method string, request proto.Message) (err error) {
	// marshal request
	pbRequest := []byte{}
	if request != nil {
		if pbRequest, err = proto.Marshal(request); err != nil {
			return
		}
	}

	// generate header
	header := &wire.RequestHeader{
		Id:     id,
		Method: method,
	}

	// check header size
	pbHeader := []byte{}
	if pbHeader, err = proto.Marshal(header); err != nil {
		return
	}
	if len(pbHeader) > 1024 {
		return fmt.Errorf("protobufc.writeRequest: header larger than max_header_len: %d", len(pbHeader))
	}

	// send header (more)
	if err = sendFrame(w, pbHeader); err != nil {
		return
	}

	// send body (end)
	return sendFrame(w, pbRequest)
}

func readRequestHeader(r io.Reader, header *wire.RequestHeader) (err error) {
	// recv header (more)
	pbHeader := []byte{}
	if pbHeader, err = recvFrame(r); err != nil {
		return
	}

	// Marshal Header
	return proto.Unmarshal(pbHeader, header)
}

func readRequestBody(r io.Reader, header *wire.RequestHeader, request proto.Message) (err error) {
	// recv body (end)
	pbRequest := []byte{}
	if pbRequest, err = recvFrame(r); err != nil {
		return
	}

	// Unmarshal to proto message
	if request == nil {
		return
	}
	return proto.Unmarshal(pbRequest, request)
}

func writeResponse(w io.Writer, id uint64, serr string, response proto.Message) (err error) {
	// clean response if error
	if serr != "" {
		response = nil
	}

	// marshal response
	pbResponse := []byte{}
	if response != nil {
		if pbResponse, err = proto.Marshal(response); err != nil {
			return
		}
	}

	// generate header
	header := &wire.ResponseHeader{
		Id:    id,
		Error: serr,
	}

	// check header size
	pbHeader := []byte{}
	if pbHeader, err = proto.Marshal(header); err != nil {
		return
	}

	// send header (more)
	if err = sendFrame(w, pbHeader); err != nil {
		return
	}

	// send body (end)
	if err = sendFrame(w, pbResponse); err != nil {
		return
	}

	return nil
}

func readResponseHeader(r io.Reader, header *wire.ResponseHeader) error {
	// recv header (more)
	pbHeader, err := recvFrame(r)
	if err != nil {
		return err
	}

	// Marshal Header
	return proto.Unmarshal(pbHeader, header)
}

func readResponseBody(r io.Reader, header *wire.ResponseHeader, response proto.Message) (err error) {
	// recv body (end)
	pbResponse := []byte{}
	if pbResponse, err = recvFrame(r); err != nil {
		return
	}

	// Unmarshal to proto message
	if response == nil {
		return
	}
	return proto.Unmarshal(pbResponse, response)
}
