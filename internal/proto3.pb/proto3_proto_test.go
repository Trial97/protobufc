// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto3_proto

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"reflect"
	"testing"
)

type tEchoService struct {
	private int
}

func (p *tEchoService) Echo(in *Message, out *Message) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(in); err != nil {
		return err
	}
	if err := gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(out); err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	go func() {
		if err := ListenAndServeEchoService("tcp", ":9527", new(tEchoService)); err != nil {
			log.Fatal(err)
		}
	}()
	os.Exit(m.Run())
}

func TestEchoService(t *testing.T) {
	c, err := DialEchoService("tcp", "127.0.0.1:9527")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	in := Message{
		Name:         "github.com/chai2010/protorpc",
		Hilarity:     Message_PUNS,
		HeightInCm:   13,
		Data:         []byte("bin data"),
		ResultCount:  2<<35 + 1,
		TrueScotsman: true,
		Score:        3.14,
		Key:          []uint64{1, 1001},
		Nested:       &Nested{Bunny: "{{Bunny}}"},
		Terrain: map[string]*Nested{
			"A": &Nested{Bunny: "{{A}}"},
			"B": &Nested{Bunny: "{{B}}"},
		},
	}

	var out Message
	if err = c.Echo(&in, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&in, &out) {
		t.Fatalf("not euqal, got = %v\n", &out)
	}
}
