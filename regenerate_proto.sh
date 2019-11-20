#!/bin/sh

protoc --gogofaster_out=. -I. -I${GOPATH}/src  -I${GOPATH}/src/github.com/gogo/protobuf/protobuf  ./examples/message.pb/*.proto

protoc --gogofaster_out=. -I. -I${GOPATH}/src  -I${GOPATH}/src/github.com/gogo/protobuf/protobuf  ./wire.pb/*.proto
