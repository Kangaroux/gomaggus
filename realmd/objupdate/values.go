package objupdate

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/phuslu/log"
)

// https://gtker.com/wow_messages/docs/updateflag.html#client-version-335
type UpdateFlag uint16

const (
	UpdateFlagNone               UpdateFlag = 0x0
	UpdateFlagSelf               UpdateFlag = 0x1
	UpdateFlagTransport          UpdateFlag = 0x2
	UpdateFlagHasAttackingTarget UpdateFlag = 0x4
	UpdateFlagLowGuid            UpdateFlag = 0x8
	UpdateFlagHighGuid           UpdateFlag = 0x10
	UpdateFlagLiving             UpdateFlag = 0x20
	UpdateFlagHasPosition        UpdateFlag = 0x40
	UpdateFlagVehicle            UpdateFlag = 0x80
	UpdateFlagPosition           UpdateFlag = 0x100
	UpdateFlagRotation           UpdateFlag = 0x200
)

type valueField struct {
	mask  FieldMask
	value []uint32
}

// Values stores the data necessary for building the values block for the object update packet.
// The values block begins with a mask that declares what fields are present, followed by the
// field values in a specific order. Field ordering is managed by Values, the caller can set
// field values in whatever order they like.
type Values struct {
	mask        ValuesMask
	fields      []*valueField
	objBuilder  *ObjectValues
	unitBuilder *UnitValues
}

// Bytes returns the complete little-endian byte array of the field mask and values.
func (v *Values) Bytes() []byte {
	buf := bytes.Buffer{}
	buf.Write(v.mask.Bytes())

	// Fields need to be sorted by offset so they are written in the correct order
	sort.Slice(v.fields, func(i, j int) bool {
		return v.fields[i].mask.Offset < v.fields[j].mask.Offset
	})

	for _, field := range v.fields {
		binary.Write(&buf, binary.LittleEndian, field.value)
		log.Trace().
			Str("field", field.mask.Name).
			Uints32("value", field.value).
			Msg("value field")
	}

	return buf.Bytes()
}

func (v *Values) Object() *ObjectValues {
	if v.objBuilder == nil {
		v.objBuilder = &ObjectValues{buf: v}
	}
	return v.objBuilder
}

func (v *Values) Unit() *UnitValues {
	if v.unitBuilder == nil {
		v.unitBuilder = &UnitValues{buf: v}
	}
	return v.unitBuilder
}

func (v *Values) addField(field *valueField) {
	// Has this field already been added?
	if v.mask.FieldMask(field.mask) {

		// Find and replace the field
		for i := range v.fields {
			if v.fields[i].mask == field.mask {
				// Setting the same field twice smells like a logic error
				log.Warn().
					Str("mask", field.mask.String()).
					Any("oldval", v.fields[i].value).
					Any("newval", v.fields[i].value).
					Msg("overwrote existing field")

				v.fields[i] = field
				return
			}
		}
	}

	v.fields = append(v.fields, field)
	v.mask.SetFieldMask(field.mask)
}
