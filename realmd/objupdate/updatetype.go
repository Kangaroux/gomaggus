package objupdate

import "github.com/kangaroux/gomaggus/realmd"

// UpdateType selects the payload to use for object update packets.
type UpdateType byte

const (
	UpdateTypePartial           UpdateType = 0
	UpdateTypeMovement          UpdateType = 1
	UpdateTypeCreateObject      UpdateType = 2
	UpdateTypeCreateNewObject   UpdateType = 3 // Unused
	UpdateTypeOutOfRangeObjects UpdateType = 4
	UpdateTypeNearObjects       UpdateType = 5 // Unused
)

// PartialUpdate contains changes to an object's non-movement properties (hp, level, etc).
type PartialUpdate struct {
	Guid   realmd.PackedGuid
	Values *Values
}

// MovementUpdate contains changes to an object's movement properties (position, speed, etc).
type MovementUpdate struct {
	Guid   realmd.PackedGuid
	Values *MovementValues
}

// CreateObject contains data for objects which have just appeared to a player. An object could appear
// if it just spawned, entered the player's range, became visible, etc.
type CreateObject struct {
	Guid     realmd.PackedGuid
	Type     ObjectType
	Movement *MovementValues
	Values   *Values
}

// ProximityUpdate contains a list of object guids that have left the player's range.
type ProximityUpdate struct {
	Count uint32
	Guids []realmd.PackedGuid
}
