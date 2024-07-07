package value

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectDataLayout(t *testing.T) {
	m := getStructLayout(reflect.ValueOf(ObjectData{}))
	assert.Equal(t, 6, m.size)
	assert.Equal(t, 5, len(m.sections))

	expected := []struct {
		blockStart int
		size       int
	}{
		{0, 2},
		{2, 1},
		{3, 1},
		{4, 1},
		{5, 1},
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
	assert.Equal(t, 1178, m.size)
	assert.Equal(t, 76, len(m.sections))

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
		{175, 1},
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
