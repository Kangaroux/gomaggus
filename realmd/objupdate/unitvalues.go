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

type UnitFlag uint32

const (
	ServerControlled    UnitFlag = 0x1
	NonAttackable       UnitFlag = 0x2
	RemoveClientControl UnitFlag = 0x4
	PlayerControlled    UnitFlag = 0x8
	Rename              UnitFlag = 0x10
	PetAbandon          UnitFlag = 0x20
	_                   UnitFlag = 0x40 // Unknown
	_                   UnitFlag = 0x80 // Unknown
	OOCNotAttackable    UnitFlag = 0x100
	Passive             UnitFlag = 0x200
	IsLooting           UnitFlag = 0x400  // Unknown
	IsPetInCombat       UnitFlag = 0x800  // Unknown
	PVP                 UnitFlag = 0x1000 // Moved to FieldMaskUnitBytes2
	IsSilenced          UnitFlag = 0x2000
	IsPersuaded         UnitFlag = 0x4000
	Swimming            UnitFlag = 0x8000
	RemoveAttackIcon    UnitFlag = 0x10000
	IsPacified          UnitFlag = 0x20000
	IsStunned           UnitFlag = 0x40000
	InCombat            UnitFlag = 0x80000
	InTaxiFlight        UnitFlag = 0x100000
	Disarmed            UnitFlag = 0x200000
	Confused            UnitFlag = 0x400000
	Fleeing             UnitFlag = 0x800000
	Possessed           UnitFlag = 0x1000000
	NotSelectable       UnitFlag = 0x2000000
	Skinnable           UnitFlag = 0x4000000
	AurasVisible        UnitFlag = 0x8000000
	_                   UnitFlag = 0x10000000 // Unknown
	NoImplicitEmotes    UnitFlag = 0x20000000
	Sheathe             UnitFlag = 0x40000000
	NoKillReward        UnitFlag = 0x80000000
)

func (v *UnitValues) Flags(flags UnitFlag) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitFlags,
		value: []uint32{uint32(flags)},
	})
}

func (v *UnitValues) Strength(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitStrength,
		value: []uint32{val},
	})
}

func (v *UnitValues) Agility(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitAgility,
		value: []uint32{val},
	})
}

func (v *UnitValues) Stamina(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitStamina,
		value: []uint32{val},
	})
}

func (v *UnitValues) Intellect(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitIntellect,
		value: []uint32{val},
	})
}

func (v *UnitValues) Spirit(val uint32) {
	v.buf.addField(&valueField{
		mask:  FieldMaskUnitSpirit,
		value: []uint32{val},
	})
}
