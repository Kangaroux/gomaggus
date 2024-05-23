package main

import (
	"bytes"
	"math/big"
)

type ByteArray struct {
	bigEndian bool
	data      []byte
}

// NewByteArray creates an endian-aware byte array from `data`. If `size` is greater than zero,
// the data will be left-padded with zeroes (see `LeftPad()`).
func NewByteArray(data []byte, size int, bigEndian bool) *ByteArray {
	ba := &ByteArray{bigEndian: bigEndian, data: data}
	if size > 0 {
		ba.LeftPad(size)
	}
	return ba
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

// LeftPad pads the left side of the byte array with zeroes if the byte array is smaller than `length`.
func (ba *ByteArray) LeftPad(length int) *ByteArray {
	dataLen := len(ba.data)

	if dataLen < length {
		padded := make([]byte, length)
		copy(padded[length-dataLen:], ba.data)
		ba.data = padded
	}

	return ba
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
func padBytes(data []byte, length int) []byte {
	dataLen := len(data)

	if dataLen >= length {
		return data
	}

	ret := make([]byte, length)
	copy(ret[length-dataLen:], data)
	return ret
}
