package objupdate

import (
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// UnitBuilder builds the values for UNIT_* types.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type UnitBuilder struct {
	buf *ValueBuffer
}

func (b *UnitBuilder) ClassRaceGenderPower(race model.Race, class model.Class, gender model.Gender, powerType realmd.PowerType) {
	val := uint32(realmd.PowerTypeForClass(class)) |
		uint32(gender)<<8 |
		uint32(class)<<16 |
		uint32(race)<<24

	b.buf.addField(&valueField{
		mask:  FieldMaskUnitBytes0,
		value: []uint32{val},
	})
}

func (b *UnitBuilder) Health(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitHealth,
		value: []uint32{val},
	})
}

func (b *UnitBuilder) MaxHealth(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitMaxHealth,
		value: []uint32{val},
	})
}

func (b *UnitBuilder) Level(val uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitLevel,
		value: []uint32{val},
	})
}

func (b *UnitBuilder) Faction(race model.Race) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitFactionTemplate,
		value: []uint32{uint32(race)},
	})
}

func (b *UnitBuilder) DisplayModel(modelId uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitDisplayId,
		value: []uint32{uint32(modelId)},
	})
}

// ???
func (b *UnitBuilder) NativeDisplayModel(modelId uint32) {
	b.buf.addField(&valueField{
		mask:  FieldMaskUnitNativeDisplayId,
		value: []uint32{uint32(modelId)},
	})
}
