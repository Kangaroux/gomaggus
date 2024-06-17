package objupdate

import (
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// UnitValueBuilder builds the values for UNIT_* types.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type UnitValueBuilder struct {
	buf *ValueBuffer
}

func (b *UnitValueBuilder) RaceClassGenderPower(race model.Race, class model.Class, gender model.Gender, powerType realmd.PowerType) {
	val := uint32(race) |
		uint32(class)<<8 |
		uint32(gender)<<16 |
		uint32(realmd.PowerTypeForClass(class))<<24

	b.buf.addField(&valueField{
		mask:  FieldMaskUnitBytes0,
		value: []uint32{val},
	})
}

func (b *UnitValueBuilder) Health(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitHealth,
		value: []uint32{val},
	})
}

func (b *UnitValueBuilder) MaxHealth(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitMaxHealth,
		value: []uint32{val},
	})
}

func (b *UnitValueBuilder) Level(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitLevel,
		value: []uint32{val},
	})
}

func (b *UnitValueBuilder) Faction(race model.Race) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitFactionTemplate,
		value: []uint32{uint32(race)},
	})
}

func (b *UnitValueBuilder) DisplayModel(modelId uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitDisplayId,
		value: []uint32{uint32(modelId)},
	})
}

// ???
func (b *UnitValueBuilder) NativeDisplayModel(modelId uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitNativeDisplayId,
		value: []uint32{uint32(modelId)},
	})
}
