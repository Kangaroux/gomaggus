package values

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectDataLayout(t *testing.T) {
	m := getStructLayout(reflect.ValueOf(ObjectData{}))
	assert.Equal(t, ObjectDataSize, m.size)
	assert.Equal(t, 4, len(m.sections))

	expected := []struct {
		blockStart int
		size       int
	}{
		{0, 2},
		{2, 1},
		{3, 1},
		{4, 1},
		// {5, 1}, padding
	}

	assert.Equal(t, len(m.sections), len(expected))

	for i := 0; i < len(expected); i++ {
		t.Run(fmt.Sprintf("block-%d", i), func(t *testing.T) {
			assert.Equal(t, expected[i].blockStart, m.sections[i].blockStart)
			assert.Equal(t, expected[i].size, m.sections[i].size)
		})
	}
}

func TestUnitDataLayout(t *testing.T) {
	m := getStructLayout(reflect.ValueOf(UnitData{}))
	assert.Equal(t, UnitDataSize, m.size)
	assert.Equal(t, 81, len(m.sections))

	expected := []struct {
		blockStart int
		size       int
	}{
		{0, 2},
		{2, 2},
		{4, 2},
		{6, 2},
		{8, 2},
		{10, 2},
		{12, 2},
		{14, 2},
		{16, 1},
		{17, 1},
		{18, 1},
		{19, 1},
		{20, 1},
		{21, 1},
		{22, 1},
		{23, 1},
		{24, 1},
		{25, 1},
		{26, 1},
		{27, 1},
		{28, 1},
		{29, 1},
		{30, 1},
		{31, 1},
		{32, 1},
		{33, 1},
		{34, 7},
		{41, 7},
		{48, 1},
		{49, 1},
		{50, 3},
		{53, 1},
		{54, 1},
		{55, 1},
		{56, 1},
		{57, 1},
		{58, 1},
		{59, 1},
		{60, 1},
		{61, 1},
		{62, 1},
		{63, 1},
		{64, 1},
		{65, 1},
		{66, 1},
		{67, 1},
		{68, 1},
		{69, 1},
		{70, 1},
		{71, 1},
		{72, 1},
		{73, 1},
		{74, 1},
		{75, 1},
		{76, 1},
		{77, 1},
		{78, 1},
		{79, 1},
		{80, 1},
		{81, 1},
		{82, 1},
		{83, 5},
		{88, 5},
		{93, 7},
		{100, 7},
		{107, 7},
		{114, 1},
		{115, 1},
		{116, 1},
		{117, 1},
		{118, 1},
		{119, 1},
		{120, 1},
		{121, 1},
		{122, 1},
		{123, 1},
		{124, 1},
		{125, 7},
		{132, 7},
		{139, 1},
		{140, 1},
		// {141, 1}, padding
	}

	assert.Equal(t, len(m.sections), len(expected))

	for i := 0; i < len(expected); i++ {
		t.Run(fmt.Sprintf("block-%d", i), func(t *testing.T) {
			assert.Equal(t, expected[i].blockStart, m.sections[i].blockStart)
			assert.Equal(t, expected[i].size, m.sections[i].size)
		})
	}
}

func TestPlayerDataLayout(t *testing.T) {
	m := getStructLayout(reflect.ValueOf(PlayerData{}))
	assert.Equal(t, PlayerDataSize, m.size)
	assert.Equal(t, 75, len(m.sections))

	expected := []struct {
		blockStart int
		size       int
	}{
		{0, 2},
		{2, 1},
		{3, 1},
		{4, 1},
		{5, 1},
		{6, 1},
		{7, 1},
		{8, 1},
		{9, 1},
		{10, 125},
		{135, 38},
		{173, 1},
		{174, 1},
		// {175, 1}, padding
		{176, 46},
		{222, 32},
		{254, 56},
		{310, 14},
		{324, 24},
		{348, 64},
		{412, 64},
		{476, 2},
		{478, 6},
		{484, 2},
		{486, 1},
		{487, 1},
		{488, 384},
		{872, 2},
		{874, 1},
		{875, 1},
		{876, 1},
		{877, 1},
		{878, 1},
		{879, 1},
		{880, 1},
		{881, 1},
		{882, 1},
		{883, 1},
		{884, 7},
		{891, 1},
		{892, 1},
		{893, 128},
		{1021, 1},
		{1022, 1},
		{1023, 7},
		{1030, 7},
		{1037, 7},
		{1044, 1},
		{1045, 1},
		{1046, 1},
		{1047, 1},
		{1048, 1},
		{1049, 1},
		{1050, 1},
		{1051, 1},
		{1052, 1},
		{1053, 12},
		{1065, 12},
		{1077, 1},
		{1078, 1},
		{1079, 1},
		{1080, 1},
		{1081, 1},
		{1082, 1},
		{1083, 25},
		{1108, 21},
		{1129, 1},
		{1130, 1},
		{1131, 1},
		{1132, 25},
		{1157, 4},
		{1161, 3},
		{1164, 6},
		{1170, 6},
		{1176, 1},
		{1177, 1},
	}

	assert.Equal(t, len(m.sections), len(expected))

	for i := 0; i < len(expected); i++ {
		t.Run(fmt.Sprintf("block-%d", i), func(t *testing.T) {
			assert.Equal(t, expected[i].blockStart, m.sections[i].blockStart)
			assert.Equal(t, expected[i].size, m.sections[i].size)
		})
	}
}

func TestFieldNotMultipleOfBlockSize(t *testing.T) {
	type foo struct {
		_ [33]bool
	}

	assert.Panics(t, func() {
		getStructLayout(reflect.ValueOf(foo{}))
	})
}
