package objupdate

import "github.com/kangaroux/gomaggus/realmd"

type UpdateType byte

const (
	UpdateTypePartial           UpdateType = 0
	UpdateTypeMovement          UpdateType = 1
	UpdateTypeCreateObject      UpdateType = 2
	UpdateTypeCreateObject2     UpdateType = 3
	UpdateTypeOutOfRangeObjects UpdateType = 4
	UpdateTypeNearObjects       UpdateType = 5 // unused
)

// PartialUpdate is used for updating non-movement properties of existing objects.
type PartialUpdate struct {
	Guid   realmd.PackedGuid
	Values *Values
}

// MovementUpdate is used for updating movement (position, speed, etc.) of existing objects.
type MovementUpdate struct {
	Guid   realmd.PackedGuid
	Values *MovementValues
}

// CreateObject is used for creating new objects and setting their properties/movement.
type CreateObject struct {
	Guid     realmd.PackedGuid
	Type     ObjectType
	Movement *MovementValues
	Values   *Values
}

// ProximityUpdate is used for notifying the client about objects which are too far away. Players can
// only see objects in a certain radius around them. Objects that leave that radius will trigger a
// ProximityUpdate so the client can hide it.
type ProximityUpdate struct {
	Count uint32
	Guids []realmd.PackedGuid
}
