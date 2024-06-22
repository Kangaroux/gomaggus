package realmd

import (
	"github.com/kangaroux/gomaggus/model"
	"github.com/phuslu/log"
)

type PowerType byte

const (
	PowerTypeMana      PowerType = 0
	PowerTypeRage      PowerType = 1
	PowerTypeFocus     PowerType = 2
	PowerTypeEnergy    PowerType = 3
	PowerTypeHappiness PowerType = 4
)

func PowerTypeForClass(c model.Class) PowerType {
	switch c {
	case model.ClassWarrior:
		return PowerTypeRage

	case model.ClassPaladin, model.ClassHunter, model.ClassPriest, model.ClassShaman, model.ClassMage, model.ClassWarlock, model.ClassDruid:
		return PowerTypeMana

	case model.ClassRogue:
		return PowerTypeEnergy

	default:
		log.Warn().Byte("class", byte(c)).Msg("invalid class")
		return PowerTypeMana
	}
}
