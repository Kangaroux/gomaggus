package srpv2

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
	for i := 0; i < n/2; i++ {
		data[i], data[n-i-1] = data[n-i-1], data[i]
	}
	return data
}

func toInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(data)
}
