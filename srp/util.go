package srp

import "math/big"

func pad(length int, data []byte) []byte {
	dataLen := len(data)
	if dataLen == length {
		return data
	}
	ret := make([]byte, length)
	copy(ret[length-dataLen:], data)
	return ret
}

func Reverse(data []byte) []byte {
	n := len(data)
	newData := make([]byte, n)
	for i := 0; i < n; i++ {
		newData[i] = data[n-i-1]
	}
	return newData
}

func toInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data)
}
