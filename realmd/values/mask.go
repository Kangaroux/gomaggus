package values

import "encoding/binary"

const (
	// The field with the largest offset is player pet spell power at 0x52D (1325)
	// with a block count of 1.
	largestBit = 1325 + 1

	// The number of mask values needed to fit the largest offset value
	maskSize = 1 + largestBit/32
)

type blockMask struct {
	mask         [maskSize]uint32
	largestIndex int
}

// Update sets the mask bits corresponding to the struct sections. The offset is added to
// the blockStart of each section.
func (m *blockMask) Update(sections []structSection, offset int) {
	var index int
	var bitIndex int

	for _, section := range sections {
		bit := (section.blockStart + offset)

		for i := 0; i < section.size; i++ {
			index = (bit + i) / 32
			bitIndex = (bit + i) % 32
			m.mask[index] |= 1 << bitIndex
		}

		if index > m.largestIndex {
			m.largestIndex = index
		}
	}
}

// Mask returns a slice of the block mask. The slice is the minimum possible length without
// any trailing zeroes.
func (m *blockMask) Mask() []uint32 {
	return m.mask[:m.largestIndex+1]
}

// Bytes returns the mask as a little endian byte array.
func (m *blockMask) Bytes() []byte {
	mask := m.Mask()
	data := make([]byte, len(mask)*4)

	for i := 0; i < len(mask); i++ {
		binary.LittleEndian.PutUint32(data[i*4:], mask[i])
	}

	return data
}
