package values

import (
	"reflect"

	"github.com/kangaroux/gomaggus/realmd"
)

const (
	ObjectDataOffset = 0
	ObjectDataSize   = 6
)

type ObjectData struct {
	guid    realmd.Guid
	objType uint32
	entry   uint32
	scaleX  float32
	_       [4]byte // padding

	dirty *dirtyValues `value:"END"`
}

func NewObjectData() *ObjectData {
	return &ObjectData{
		dirty: newDirtyValues(getStructLayout(reflect.ValueOf(ObjectData{}))),
	}
}

func (o *ObjectData) Marshal(onlyDirty bool) ([]byte, []structSection) {
	return marshalValues(o, onlyDirty, o.dirty)
}

func (o *ObjectData) GUID() realmd.Guid {
	return o.guid
}

func (o *ObjectData) SetGUID(val realmd.Guid) {
	o.guid = val
	o.dirty.Flag("guid")
}

func (o *ObjectData) Type() uint32 {
	return o.objType
}

func (o *ObjectData) SetType(val uint32) {
	o.objType = val
	o.dirty.Flag("objType")
}

func (o *ObjectData) Entry() uint32 {
	return o.entry
}

func (o *ObjectData) SetEntry(val uint32) {
	o.entry = val
	o.dirty.Flag("entry")
}

func (o *ObjectData) ScaleX() float32 {
	return o.scaleX
}

func (o *ObjectData) SetScaleX(val float32) {
	o.scaleX = val
	o.dirty.Flag("scaleX")
}
