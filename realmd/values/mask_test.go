package values

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	testCases := []struct {
		sections []structSection
		offset   int
		expected []uint32
	}{
		{
			sections: []structSection{},
			offset:   0,
			expected: []uint32(nil),
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 1},
			},
			offset:   0,
			expected: []uint32{1},
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 2},
			},
			offset:   0,
			expected: []uint32{3},
		},
		{
			sections: []structSection{
				{blockStart: 32, size: 1},
			},
			offset:   0,
			expected: []uint32{0, 1},
		},
		{
			sections: []structSection{
				{blockStart: 31, size: 2},
			},
			offset:   0,
			expected: []uint32{0x80000000, 1},
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 1},
			},
			offset:   1, // offset > 0
			expected: []uint32{2},
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 1}, {blockStart: 1, size: 1},
			},
			offset:   0,
			expected: []uint32{3},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			m := blockMask{}
			m.Update(testCase.sections, testCase.offset)
			assert.Equal(t, testCase.expected, m.Mask())
		})
	}
}

func TestMaskMaxBlock(t *testing.T) {
	m := blockMask{}

	assert.NotPanics(t, func() {
		m.Update([]structSection{{blockStart: maskSize*32 - 1, size: 1}}, 0)
	})

	assert.Panics(t, func() {
		m.Update([]structSection{{blockStart: maskSize * 32, size: 1}}, 0)
	})
}

func TestMaskBytes(t *testing.T) {
	testCases := []struct {
		sections []structSection
		expected []byte
	}{
		{
			sections: []structSection{},
			expected: []byte{0},
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 1},
			},
			expected: []byte{1, 1, 0, 0, 0},
		},
		{
			sections: []structSection{
				{blockStart: 0, size: 2},
			},
			expected: []byte{1, 3, 0, 0, 0},
		},
		{
			sections: []structSection{
				{blockStart: 32, size: 1},
			},
			expected: []byte{2, 0, 0, 0, 0, 1, 0, 0, 0},
		},
		{
			sections: []structSection{
				{blockStart: 31, size: 2},
			},
			expected: []byte{2, 0, 0, 0, 0x80, 1, 0, 0, 0},
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("case-%d", i), func(t *testing.T) {
			m := blockMask{}
			m.Update(testCase.sections, 0)
			assert.Equal(t, testCase.expected, m.Bytes())
		})
	}
}
