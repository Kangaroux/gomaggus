package objupdate

import (
	"math"

	"github.com/kangaroux/gomaggus/realmd"
)

type ObjectType byte

const (
	ObjectTypeObject        ObjectType = 0
	ObjectTypeItem          ObjectType = 1
	ObjectTypeContainer     ObjectType = 2
	ObjectTypeUnit          ObjectType = 3
	ObjectTypePlayer        ObjectType = 4
	ObjectTypeGameObject    ObjectType = 5
	ObjectTypeDynamicObject ObjectType = 6
	ObjectTypeCorpse        ObjectType = 7
)

// ObjectValues provides an interface for setting values for objects. Values can be added in any order.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type ObjectValues struct {
	buf *Values
}

func (v *ObjectValues) Guid(guid realmd.Guid) {
	v.buf.addField(&valueField{
		mask:  FieldMaskObjectGuid,
		value: []uint32{uint32(guid), uint32(guid >> 32)},
	})
}

func (v *ObjectValues) Type(types ...ObjectType) {
	val := uint32(0)

	// Convert types to bitmask
	for _, t := range types {
		val |= 1 << t
	}

	v.buf.addField(&valueField{
		mask:  FieldMaskObjectType,
		value: []uint32{val},
	})
}

func (v *ObjectValues) ScaleX(val float32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskObjectScaleX,
		value: []uint32{math.Float32bits(val)},
	})
}
