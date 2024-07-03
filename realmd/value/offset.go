package value

import (
	"fmt"
	"reflect"
)

const (
	tagName   = "value"
	endMarker = "END"
)

// fieldOffsets32 inspects the struct v, returning a list of field offsets and the effective
// size of v. The offsets and effective size count the number of 32 bit chunks. An offset
// or size of 1 means (1) 32 bit chunk, or (4) bytes. Fields named "_" are skipped. If a
// field has a struct tag of `value:"END"`, fieldOffsets32 will immediately return. This can
// be used to store metadata at the end of the struct.
func fieldOffsets32(v any) ([]uint16, int) {
	var offsets []uint16

	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		panic("v must be a struct")
	}

	size := 0
	prevOffset := -1

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.Tag.Get(tagName) == endMarker {
			break
		}

		// It's important that field.Offset is not used since it takes into account memory alignment.
		// We want the field's offset in a byte array, not in memory.
		offset := size / 4
		size += int(field.Type.Size())

		if field.Name == "_" {
			continue
		}

		// Add the offset if it's unique. The offset is rounded down if it's not an even multiple of 4.
		// This allows 8 and 16 bit fields to be separated in the struct but represented as a single
		// packed value once encoded.
		if offset != prevOffset {
			offsets = append(offsets, uint16(offset))
			prevOffset = offset
		}
	}

	if size%4 != 0 {
		panic(fmt.Sprintf("size of %s should be an even multiple of 4 (got %d)", t.Name(), size))
	}

	return offsets, int(size) / 4
}
