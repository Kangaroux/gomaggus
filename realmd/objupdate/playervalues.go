package objupdate

// PlayerValues provides an interface for setting values for players. Values can be added in any order.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type PlayerValues struct {
	buf *Values
}

const (
	skinMask      = 0xFF
	faceMask      = 0xFF00
	hairStyleMask = 0xFF0000
	hairColorMask = 0xFF000000
)

func (v *PlayerValues) Skin(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerFieldBytes,
		value:   []uint32{uint32(val)},
		bitmask: []uint32{skinMask},
	})
}

func (v *PlayerValues) Face(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerFieldBytes,
		value:   []uint32{uint32(val) << 8},
		bitmask: []uint32{faceMask},
	})
}
func (v *PlayerValues) HairStyle(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerFieldBytes,
		value:   []uint32{uint32(val) << 16},
		bitmask: []uint32{hairStyleMask},
	})
}
func (v *PlayerValues) HairColor(val uint8) {
	v.buf.addField(&valueField{
		mask:    FieldMaskPlayerFieldBytes,
		value:   []uint32{uint32(val) << 24},
		bitmask: []uint32{hairColorMask},
	})
}
