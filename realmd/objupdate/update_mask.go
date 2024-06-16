package realmd

type UpdateMask struct {
	largestBit int
	mask       []uint32
}

func NewUpdateMask() *UpdateMask {
	return &UpdateMask{mask: make([]uint32, 16)}
}

// Mask returns the smallest []uint32 to represent all of the mask bits that were set.
func (m *UpdateMask) Mask() []uint32 {
	largestBitIndex := m.largestBit / 32
	return m.mask[:largestBitIndex+1]
}

// SetFieldMask sets all the bits necessary for the provided field mask.
func (m *UpdateMask) SetFieldMask(fieldMask FieldMask) {
	for i := 0; i < fieldMask.Size; i++ {
		m.SetBit(int(fieldMask.Offset) + i)
	}
}

// SetBit sets the nth bit in the update mask. The bit is zero-indexed with the first bit being zero.
func (m *UpdateMask) SetBit(bit int) {
	index := bit / 32
	bitPos := bit % 32
	m.resize(index)

	if bit > m.largestBit {
		m.largestBit = bit
	}

	m.mask[index] |= 1 << bitPos
}

// Resizes the mask to fit up to n uint32s.
func (m *UpdateMask) resize(n int) {
	if len(m.mask) > n {
		return
	}

	// Grow the array exponentially
	newSize := len(m.mask)
	newSize *= newSize

	// If it's still too small just use the desired size
	if newSize < n {
		newSize = n
	}

	oldMask := m.mask
	m.mask = make([]uint32, newSize)
	copy(m.mask, oldMask)
}
