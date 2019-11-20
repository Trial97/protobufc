# Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

PROTO_FILES=$(sort $(wildcard ./*.proto))

default: $(PROTO_FILES) Makefile
	go install github.com/cgrates/protobufc/protoc-gen-protorpc
	protoc --go_out=. ${PROTO_FILES}
	protoc --protorpc_out=. ${PROTO_FILES}
	go test

clean:
