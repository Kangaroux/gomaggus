package value

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
			expected: []uint32{0},
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
