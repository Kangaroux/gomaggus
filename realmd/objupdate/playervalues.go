package objupdate

// PlayerValues provides an interface for setting values for players. Values can be added in any order.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type PlayerValues struct {
	buf *Values
}

const (
	skinColorMask = 0xFF
	faceMask      = 0xFF00
	hairStyleMask = 0xFF0000
	hairColorMask = 0xFF000000
)

func (v *PlayerValues) SkinColor(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes1,
		value:   []uint32{uint32(val)},
		bitmask: []uint32{skinColorMask},
	})
}

func (v *PlayerValues) Face(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes1,
		value:   []uint32{uint32(val) << 8},
		bitmask: []uint32{faceMask},
	})
}

func (v *PlayerValues) HairStyle(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes1,
		value:   []uint32{uint32(val) << 16},
		bitmask: []uint32{hairStyleMask},
	})
}

func (v *PlayerValues) HairColor(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes1,
		value:   []uint32{uint32(val) << 24},
		bitmask: []uint32{hairColorMask},
	})
}

const (
	extraCosmeticMask    = 0xFF
	bankBagSlotCountMask = 0xFF00
	restStateMask        = 0xFF0000
	playerGenderMask     = 0xFF000000 // different from unit gender?
)

func (v *PlayerValues) ExtraCosmetic(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes2,
		value:   []uint32{uint32(val)},
		bitmask: []uint32{extraCosmeticMask},
	})
}

func (v *PlayerValues) BankBagSlotCount(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes2,
		value:   []uint32{uint32(val) << 8},
		bitmask: []uint32{bankBagSlotCountMask},
	})
}

func (v *PlayerValues) RestState(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes2,
		value:   []uint32{uint32(val) << 16},
		bitmask: []uint32{restStateMask},
	})
}

func (v *PlayerValues) PlayerGender(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerBytes2,
		value:   []uint32{uint32(val) << 24},
		bitmask: []uint32{playerGenderMask},
	})
}
