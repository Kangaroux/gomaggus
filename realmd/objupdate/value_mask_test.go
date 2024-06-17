package objupdate_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	o "github.com/kangaroux/gomaggus/realmd/objupdate"
	"github.com/stretchr/testify/assert"
)

func TestValueMask(t *testing.T) {
	cases := []struct {
		masks    []o.FieldMask
		expected []uint32
	}{
		{
			masks:    []o.FieldMask{},
			expected: []uint32{},
		},
		{
			masks:    []o.FieldMask{{Size: 1, Offset: 0}},
			expected: []uint32{0x1},
		},
		{
			masks:    []o.FieldMask{{Size: 32, Offset: 0}},
			expected: []uint32{0xFFFFFFFF},
		},
		{
			masks:    []o.FieldMask{{Size: 33, Offset: 0}},
			expected: []uint32{0xFFFFFFFF, 0x1},
		},
		{
			masks:    []o.FieldMask{{Size: 1, Offset: 32}},
			expected: []uint32{0x0, 0x1},
		},
		{
			masks:    []o.FieldMask{{Size: 2, Offset: 31}},
			expected: []uint32{0x80000000, 0x1},
		},
		{
			masks:    []o.FieldMask{{Size: 1, Offset: 0}, {Size: 2, Offset: 2}},
			expected: []uint32{0xD},
		},
		{
			// Largest field offset
			masks: []o.FieldMask{o.FieldMaskPlayerPetSpellPower},
			expected: []uint32{
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x2000,
			},
		},
	}

	buf := bytes.Buffer{}

	for _, c := range cases {
		vm := o.ValueMask{}

		for _, fm := range c.masks {
			vm.SetFieldMask(fm)
		}

		buf.WriteByte(byte(len(c.expected)))                // mask size
		binary.Write(&buf, binary.LittleEndian, c.expected) // convert BE []uint32 to LE []byte

		assert.Equal(t, buf.Bytes(), vm.Bytes())

		buf.Reset()
	}
}
