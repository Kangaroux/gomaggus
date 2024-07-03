package objupdate

import (
	"testing"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/mixcode/binarystruct"
	"github.com/stretchr/testify/assert"
)

func TestRaceClassGenderPower(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.Race(4)
	v.Class(3)
	v.Gender(2)
	v.PowerType(1)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitRaceClassGenderPower),
			Values: []uint32{0x01020304},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestHealth(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.Health(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitHealth),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestMaxHealth(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.MaxHealth(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitMaxHealth),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestLevel(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.Level(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitLevel),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestFaction(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.Faction(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitFactionTemplate),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestDisplayModel(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.DisplayModel(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitDisplayId),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}

func TestNativeDisplayModel(t *testing.T) {
	v := UnitValues{buf: &Values{}}
	v.NativeDisplayModel(123)
	expected := internal.MustMarshal(
		&valueBlock{
			Mask:   makeMask(FieldMaskUnitNativeDisplayId),
			Values: []uint32{123},
		},
		binarystruct.LittleEndian,
	)

	assert.Equal(t, expected, v.buf.Bytes())
}
