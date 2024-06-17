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

// ObjectBuilder builds the values for OBJECT_* types.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type ObjectBuilder struct {
	buf *ValueBuffer
}

func (b *ObjectBuilder) Guid(guid realmd.Guid) {
	b.buf.addField(&valueField{
		mask:  FieldMaskObjectGuid,
		value: []uint32{uint32(guid >> 32), uint32(guid)},
	})
}

func (b *ObjectBuilder) Type(t ObjectType) {
	b.buf.addField(&valueField{
		mask:  FieldMaskObjectType,
		value: []uint32{uint32(t)},
	})
}

func (b *ObjectBuilder) ScaleX(val float32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskObjectScaleX,
		value: []uint32{math.Float32bits(val)},
	})
}
