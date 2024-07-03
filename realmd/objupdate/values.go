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
	mask    FieldMask
	value   []uint32
	bitmask []uint32
}

// Values stores the data necessary for building the values block for the object update packet.
// The values block begins with a mask that declares what fields are present, followed by the
// field values in a specific order. Field ordering is managed by Values, the caller can set
// field values in whatever order they like.
type Values struct {
	mask   ValuesMask
	fields []*valueField
	object *ObjectValues
	player *PlayerValues
	unit   *UnitValues
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
	if v.object == nil {
		v.object = &ObjectValues{buf: v}
	}
	return v.object
}

func (v *Values) Player() *PlayerValues {
	if v.player == nil {
		v.player = &PlayerValues{buf: v}
	}
	return v.player
}

func (v *Values) Unit() *UnitValues {
	if v.unit == nil {
		v.unit = &UnitValues{buf: v}
	}
	return v.unit
}

func (v *Values) addField(field *valueField) {
	// Field has not been added yet
	if !v.mask.FieldMask(field.mask) {
		v.fields = append(v.fields, field)
		v.mask.SetFieldMask(field.mask)
		return
	}

	// Find existing field and update it
	for i := range v.fields {
		existing := v.fields[i]
		if existing.mask != field.mask {
			continue
		}

		// If a bitmask was provided, overwrite only the masked bits
		if len(field.bitmask) > 0 {

			// Expand the existing value to fit the provided bitmask
			if len(field.bitmask) > len(existing.value) {
				for j := len(existing.value); j < len(field.bitmask); j++ {
					existing.value = append(existing.value, 0)
				}
			}

			for j := range existing.bitmask {
				// Zero the masked bits first in case there was a value there previously
				existing.value[j] &= ^field.bitmask[j]
				existing.value[j] |= field.value[j] & field.bitmask[j]
			}
		} else {
			// If no bitmask was provided, overwrite the whole field
			v.fields[i].value = field.value
		}

		break
	}
}
