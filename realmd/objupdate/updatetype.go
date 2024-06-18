package objupdate

import "github.com/kangaroux/gomaggus/realmd"

type UpdateType byte

const (
	UpdateTypePartial           UpdateType = 0
	UpdateTypeMovement          UpdateType = 1
	UpdateTypeCreateObject      UpdateType = 2
	UpdateTypeCreateObject2     UpdateType = 3
	UpdateTypeOutOfRangeObjects UpdateType = 4
	UpdateTypeNearObjects       UpdateType = 5
)

type PartialUpdate struct {
	Guid   realmd.PackedGuid
	Values Values
}
