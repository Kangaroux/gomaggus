package value

import (
	"reflect"

	"github.com/kangaroux/gomaggus/realmd"
)

type Object struct {
	guid    realmd.Guid
	objType uint32
	entry   uint32
	scaleX  float32
	_       [4]byte // padding

	dirty *dirtyValues `value:"END"`
}

func NewObject() *Object {
	return &Object{
		dirty: newDirtyValues(getStructLayout(reflect.ValueOf(Object{}))),
	}
}

func (o *Object) GUID() realmd.Guid {
	return o.guid
}

func (o *Object) SetGUID(val realmd.Guid) {
	o.guid = val
	o.dirty.Flag("guid")
}

func (o *Object) Type() uint32 {
	return o.objType
}

func (o *Object) SetType(val uint32) {
	o.objType = val
	o.dirty.Flag("objType")
}

func (o *Object) Entry() uint32 {
	return o.entry
}

func (o *Object) SetEntry(val uint32) {
	o.entry = val
	o.dirty.Flag("entry")
}

func (o *Object) ScaleX() float32 {
	return o.scaleX
}

func (o *Object) SetScaleX(val float32) {
	o.scaleX = val
	o.dirty.Flag("scaleX")
}
