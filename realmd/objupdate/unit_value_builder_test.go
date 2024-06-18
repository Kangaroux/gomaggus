package objupdate

import (
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/mixcode/binarystruct"
	"github.com/stretchr/testify/assert"
)

func TestRaceClassGenderPower(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.RaceClassGenderPower(4, 3, 2, 1)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitRaceClassGenderPower),
			Values: []uint32{0x01020304},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestHealth(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.Health(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitHealth),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestMaxHealth(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.MaxHealth(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitMaxHealth),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestLevel(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.Level(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitLevel),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestFaction(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.Faction(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitFactionTemplate),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestDisplayModel(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.DisplayModel(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitDisplayId),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}

func TestNativeDisplayModel(t *testing.T) {
	b := UnitValueBuilder{buf: &ValueBuffer{}}
	b.NativeDisplayModel(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitNativeDisplayId),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, b.buf.Bytes())
}
