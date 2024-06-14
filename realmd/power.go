package realmd

import (
	"log"

	"github.com/kangaroux/gomaggus/model"
)

func PowerTypeForClass(c model.Class) PowerType {
	switch c {
	case model.ClassWarrior:
		return PowerTypeRage

	case model.ClassPaladin,
		model.ClassHunter,
		model.ClassPriest,
		model.ClassShaman,
		model.ClassMage,
		model.ClassWarlock,
		model.ClassDruid:
		return PowerTypeMana

	case model.ClassRogue:
		return PowerTypeEnergy

	default:
		log.Println("getPowerTypeForClass: got unexpected class", c)
		return PowerTypeMana
	}
}
