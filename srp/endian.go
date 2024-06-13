package srp

import (
	"math/big"

	"github.com/kangaroux/gomaggus/internal"
)

// BytesToInt returns a little endian big integer from a big endian byte array.
func BytesToInt(data []byte) *big.Int {
	return big.NewInt(0).SetBytes(internal.Reverse(data))
}

// IntToBytes returns a big endian byte array from a little endian big integer.
func IntToBytes(padding int, bi *big.Int) []byte {
	return internal.Reverse(internal.Pad(padding, bi.Bytes()))
}
