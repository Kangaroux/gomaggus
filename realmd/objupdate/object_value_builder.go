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

// ObjectValueBuilder builds the values for OBJECT_* types.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type ObjectValueBuilder struct {
	buf *ValueBuffer
}

func (b *ObjectValueBuilder) Guid(guid realmd.Guid) {
	b.buf.addField(&valueField{
		mask:  FieldMaskObjectGuid,
		value: []uint32{uint32(guid >> 32), uint32(guid)},
	})
}

func (b *ObjectValueBuilder) Type(types ...ObjectType) {
	val := uint32(0)

	// Convert types to bitmask
	for _, t := range types {
		val |= 1 << t
	}

	b.buf.addField(&valueField{
		mask:  FieldMaskObjectType,
		value: []uint32{val},
	})
}

func (b *ObjectValueBuilder) ScaleX(val float32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskObjectScaleX,
		value: []uint32{math.Float32bits(val)},
	})
}
