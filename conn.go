// Copyright 2013 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protobufc

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

func sendFrame(w io.Writer, data []byte) (err error) {
	// Allocate enough space for the biggest uvarint
	var size [binary.MaxVarintLen64]byte

	if data == nil || len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n], false); err != nil {
			return
		}
		return
	}

	// Write the size and data
	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err = write(w, size[:n], false); err != nil {
		return
	}
	return write(w, data, false)
}

func recvFrame(r io.Reader) (data []byte, err error) {
	size, err := readUvarint(r)
	if err != nil {
		return nil, err
	}
	if size == 0 {
		return
	}
	data = make([]byte, size)
	if err = read(r, data); err != nil {
		return nil, err
	}
	return
}

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func readUvarint(r io.Reader) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; ; i++ {
		var b byte
		b, err := readByte(r)
		if err != nil {
			return 0, err
		}
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return x, errors.New("protobufc: varint overflows a 64-bit integer")
			}
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
}

func write(w io.Writer, data []byte, onePacket bool) (err error) {
	if onePacket {
		_, err = w.Write(data)
		return
	}
	for index := 0; index < len(data); {
		var n int
		if n, err = w.Write(data[index:]); err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return
			}
		}
		index += n
	}
	return
}

func read(r io.Reader, data []byte) (err error) {
	for index := 0; index < len(data); {
		var n int
		if n, err = r.Read(data[index:]); err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return
			}
		}
		index += n
	}
	return
}

func readByte(r io.Reader) (c byte, err error) {
	data := make([]byte, 1)
	if err = read(r, data); err != nil {
		return
	}
	c = data[0]
	return
}
