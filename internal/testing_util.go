package internal

import (
	"encoding/csv"
	"encoding/hex"
	"os"
)

// LoadTestData reads a CSV containing test inputs and returns a 2D array of the rows and columns.
// This should only be used inside tests. Panics if an error occurs.
func LoadTestData(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	return rows
}

// DecodeHex returns a byte array parsed from the given hex string. This should only be used inside
// tests. Panics if an error occurs.
func DecodeHex(s string) []byte {
	val, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return val
}
