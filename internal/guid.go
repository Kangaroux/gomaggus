package internal

// internal.PackGuid returns a packed *little-endian* representation of an 8-byte integer. The packing works
// by creating a bit mask to mark which bytes are non-zero. Any bytes which are zero are discarded.
// The result is a byte array with the first byte as the bitmask, followed by the remaining
// undiscarded bytes. The bytes after the bitmask are little-endian.
func PackGuid(val uint64) []byte {
	// At its largest, a packed guid takes up 9 bytes (1 byte mask + 8 bytes)
	result := make([]byte, 9)
	n := 0

	for i := 0; i < 8; i++ {
		if byte(val) > 0 {
			// Set the mask bit
			result[0] |= 1 << i
			// Add the byte to the result. The loop traverses the bytes from right-to-left but they
			// are written to the result from left-to-right, which swaps it to little-endian.
			result[1] = byte(val)
			n++
		}
		// Move to the next byte
		val >>= 8
	}

	return result[:n+1]
}
