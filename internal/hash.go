package internal

// FastHash returns an *insecure* 64-bit hash value. This uses the djb2 algorithm.
// Source: http://www.cse.yorku.ca/~oz/hash.html
func FastHash(s string) int64 {
	h := int64(5381)
	for i := 0; i < len(s); i++ {
		h = ((h << 5) + h) + int64(s[i])
	}
	return h
}
