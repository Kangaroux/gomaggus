package objupdate

import (
	"bytes"
	"encoding/binary"
	"sort"
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

type ValuesBuffer struct {
	mask       UpdateMask
	fields     []*valueField
	objBuilder *ObjectBuilder
}

func (vb *ValuesBuffer) Bytes() []byte {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.LittleEndian, vb.mask.Mask())

	// Fields need to be sorted by offset so they are written in the correct order
	sort.Slice(vb.fields, func(i, j int) bool {
		return vb.fields[i].mask.Offset-vb.fields[j].mask.Offset > 0
	})

	for _, field := range vb.fields {
		binary.Write(&buf, binary.LittleEndian, field.value)
	}

	return buf.Bytes()
}

func (vb *ValuesBuffer) Objects() *ObjectBuilder {
	if vb.objBuilder == nil {
		vb.objBuilder = &ObjectBuilder{buf: vb}
	}
	return vb.objBuilder
}

func (vb *ValuesBuffer) addField(field *valueField) {
	vb.fields = append(vb.fields, field)
	vb.mask.SetFieldMask(field.mask)
}
