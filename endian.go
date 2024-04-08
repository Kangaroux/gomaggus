package main

import (
	"bytes"
	"encoding/hex"
)

type EndianBytes struct {
	data      []byte
	bigEndian bool
}

func (eb *EndianBytes) Bytes() []byte {
	return eb.data
}

func (eb *EndianBytes) SetBytes(data []byte, bigEndian bool) {
	eb.data = bytes.Clone(data)
	eb.bigEndian = bigEndian
}

func (eb *EndianBytes) ToBigEndian() {
	if eb.bigEndian {
		return
	}

	eb.swapEndian()
}

func (eb *EndianBytes) ToLittleEndian() {
	if !eb.bigEndian {
		return
	}

	eb.swapEndian()
}

func (eb *EndianBytes) swapEndian() {
	eb.bigEndian = !eb.bigEndian
	reverseBytes(eb.data, len(eb.data))
}

func BytesFromHex(s string, bigEndian bool) (*EndianBytes, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	eb := EndianBytes{}
	eb.SetBytes(b, bigEndian)
	return &eb, nil
}

func reverseBytes(data []byte, n int) {
	for i := 0; i < n/2; i++ {
		data[i], data[n-i-1] = data[n-i-1], data[i]
	}
}
