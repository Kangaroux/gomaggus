package main

import (
	"bytes"
	"math/big"
)

type ByteArray struct {
	bigEndian bool
	data      []byte
}

func NewByteArray(data []byte, bigEndian bool) *ByteArray {
	return &ByteArray{bigEndian: bigEndian, data: data}
}

func (ba *ByteArray) Bytes() []byte {
	return ba.data
}

func (ba *ByteArray) Clone() *ByteArray {
	return &ByteArray{
		bigEndian: ba.bigEndian,
		data:      bytes.Clone(ba.data),
	}
}

func (ba *ByteArray) BigEndian() *ByteArray {
	if !ba.bigEndian {
		ba = ba.Clone()
		ba.swapEndian()
	}
	return ba
}

func (ba *ByteArray) LittleEndian() *ByteArray {
	if ba.bigEndian {
		ba = ba.Clone()
		ba.swapEndian()
	}
	return ba
}

func (ba *ByteArray) BigInt() BigInteger {
	return big.NewInt(0).SetBytes(ba.data)
}

func (ba *ByteArray) swapEndian() {
	ba.bigEndian = !ba.bigEndian
	reverseBytesNoCopy(ba.data)
}

// reverses the byte array in place
func reverseBytesNoCopy(data []byte) {
	n := len(data)
	for i := 0; i < n/2; i++ {
		data[i], data[n-i-1] = data[n-i-1], data[i]
	}
}

// PadBytes pads the left side of the `data` with zeroes if `len(data)` is less than `length`.
// Returns a new byte array if padding was added, otherwise returns the original `data`.
func PadBytes(data []byte, length int) []byte {
	dataLen := len(data)

	if dataLen >= length {
		return data
	}

	ret := make([]byte, length)
	copy(ret[length-dataLen:], data)
	return ret
}
