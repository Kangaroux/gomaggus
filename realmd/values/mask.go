package values

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
