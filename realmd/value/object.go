package value

import (
	"github.com/kangaroux/gomaggus/realmd"
)

type objectField int

const (
	objGuid objectField = 1 << iota
	objType
	objEntry
	objScaleX
)

type Object struct {
	guid    realmd.Guid
	objType uint32
	entry   uint32
	scaleX  float32
	_       [4]byte // padding

	dirty objectField `value:"END"`
}

func (o *Object) GUID() realmd.Guid {
	return o.guid
}

func (o *Object) SetGUID(val realmd.Guid) {
	o.guid = val
	o.dirty |= objGuid
}

func (o *Object) Type() uint32 {
	return o.objType
}

func (o *Object) SetType(val uint32) {
	o.objType = val
	o.dirty |= objType
}

func (o *Object) Entry() uint32 {
	return o.entry
}

func (o *Object) SetEntry(val uint32) {
	o.entry = val
	o.dirty |= objEntry
}

func (o *Object) ScaleX() float32 {
	return o.scaleX
}

func (o *Object) SetScaleX(val float32) {
	o.scaleX = val
	o.dirty |= objScaleX
}
