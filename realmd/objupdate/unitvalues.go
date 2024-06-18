package objupdate

import (
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// UnitValues provides an interface for setting values for units. Values can be added in any order.
// https://gtker.com/wow_messages/types/update-mask.html#version-335
type UnitValues struct {
	buf *Values
}

func (v *UnitValues) RaceClassGenderPower(race model.Race, class model.Class, gender model.Gender, powerType realmd.PowerType) {
	val := uint32(race) |
		uint32(class)<<8 |
		uint32(gender)<<16 |
		uint32(powerType)<<24

	v.buf.addField(&valueField{
		mask:  FieldMaskUnitRaceClassGenderPower,
		value: []uint32{val},
	})
}

func (v *UnitValues) Health(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitHealth,
		value: []uint32{val},
	})
}

func (v *UnitValues) MaxHealth(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitMaxHealth,
		value: []uint32{val},
	})
}

func (v *UnitValues) Level(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitLevel,
		value: []uint32{val},
	})
}

func (v *UnitValues) Faction(race model.Race) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitFactionTemplate,
		value: []uint32{uint32(race)},
	})
}

func (v *UnitValues) DisplayModel(modelId uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitDisplayId,
		value: []uint32{uint32(modelId)},
	})
}

// ???
func (v *UnitValues) NativeDisplayModel(modelId uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitNativeDisplayId,
		value: []uint32{uint32(modelId)},
	})
}
