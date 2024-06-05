package internal

// Pad adds zero-bytes as padding to a byte array if the array is smaller than the given length.
// If the array doesn't need padding, the original array is returned, otherwise a copy with padding
// included is returned.
func Pad(length int, data []byte) []byte {
	dataLen := len(data)
	if dataLen >= length {
		return data
	}
	ret := make([]byte, length)
	copy(ret[length-dataLen:], data)
	return ret
}

// Reverse returns a copy of the given byte array in reverse order.
func Reverse(data []byte) []byte {
	n := len(data)
	newData := make([]byte, n)
	for i := 0; i < n; i++ {
		newData[i] = data[n-i-1]
	}
	return newData
}

// FastHash returns an *insecure* 64-bit hash value. This uses the djb2 algorithm.
// Source: http://www.cse.yorku.ca/~oz/hash.html
func FastHash(s string) int64 {
	h := int64(5381)
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + int64(s[i])
	}
	return h
}
