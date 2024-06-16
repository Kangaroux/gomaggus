package objupdate

type UpdateType byte

const (
	UpdateTypePartial           UpdateType = 0
	UpdateTypeMovement          UpdateType = 1
	UpdateTypeCreateObject      UpdateType = 2
	UpdateTypeCreateObject2     UpdateType = 3
	UpdateTypeOutOfRangeObjects UpdateType = 4
	UpdateTypeNearObjects       UpdateType = 5
)

// https://gtker.com/wow_messages/docs/updateflag.html#client-version-335
type UpdateFlag uint16

const (
	UpdateFlagNone               UpdateFlag = 0x0
	UpdateFlagSelf               UpdateFlag = 0x1
	UpdateFlagTransport          UpdateFlag = 0x2
	UpdateFlagHasAttackingTarget UpdateFlag = 0x4
	UpdateFlagLowGuid            UpdateFlag = 0x8
	UpdateFlagHighGuid           UpdateFlag = 0x10
	UpdateFlagLiving             UpdateFlag = 0x20
	UpdateFlagHasPosition        UpdateFlag = 0x40
	UpdateFlagVehicle            UpdateFlag = 0x80
	UpdateFlagPosition           UpdateFlag = 0x100
	UpdateFlagRotation           UpdateFlag = 0x200
)

type Builder struct {
	mask   UpdateMask
	values []uint32
}
