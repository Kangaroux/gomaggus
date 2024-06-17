package objupdate

import (
	"bytes"
	"encoding/binary"

	"github.com/kangaroux/gomaggus/realmd"
)

type MovementUpdateFlag uint16

const (
	MovementUpdateNone               MovementUpdateFlag = 0x0
	MovementUpdateSelf               MovementUpdateFlag = 0x1
	MovementUpdateTransport          MovementUpdateFlag = 0x2
	MovementUpdateHasAttackingTarget MovementUpdateFlag = 0x4
	MovementUpdateLowGuid            MovementUpdateFlag = 0x8
	MovementUpdateHighGuid           MovementUpdateFlag = 0x10
	MovementUpdateLiving             MovementUpdateFlag = 0x20
	MovementUpdateHasPosition        MovementUpdateFlag = 0x40
	MovementUpdateVehicle            MovementUpdateFlag = 0x80
	MovementUpdatePosition           MovementUpdateFlag = 0x100
	MovementUpdateRotation           MovementUpdateFlag = 0x200
)

// Encoded as 48 bits
type MovementFlag uint64

const (
	MovementFlagNone                          MovementFlag = 0x0
	MovementFlagForward                       MovementFlag = 0x1
	MovementFlagBackward                      MovementFlag = 0x2
	MovementFlagStrafeLeft                    MovementFlag = 0x4
	MovementFlagStrafeRight                   MovementFlag = 0x8
	MovementFlagLeft                          MovementFlag = 0x10
	MovementFlagRight                         MovementFlag = 0x20
	MovementFlagPitchUp                       MovementFlag = 0x40
	MovementFlagPitchDown                     MovementFlag = 0x80
	MovementFlagWalking                       MovementFlag = 0x100
	MovementFlagOnTransport                   MovementFlag = 0x200
	MovementFlagDisableGravity                MovementFlag = 0x400
	MovementFlagRoot                          MovementFlag = 0x800
	MovementFlagFalling                       MovementFlag = 0x1000
	MovementFlagFallingFar                    MovementFlag = 0x2000
	MovementFlagPendingStop                   MovementFlag = 0x4000
	MovementFlagPendingStrafeStop             MovementFlag = 0x8000
	MovementFlagPendingForward                MovementFlag = 0x10000
	MovementFlagPendingBackward               MovementFlag = 0x20000
	MovementFlagPendingStrafeLeft             MovementFlag = 0x40000
	MovementFlagPendingStrafeRight            MovementFlag = 0x80000
	MovementFlagPendingRoot                   MovementFlag = 0x100000
	MovementFlagSwimming                      MovementFlag = 0x200000
	MovementFlagAscending                     MovementFlag = 0x400000
	MovementFlagDescending                    MovementFlag = 0x800000
	MovementFlagCanFly                        MovementFlag = 0x1000000
	MovementFlagFlying                        MovementFlag = 0x2000000
	MovementFlagSplineElevation               MovementFlag = 0x4000000
	MovementFlagSplineEnabled                 MovementFlag = 0x8000000
	MovementFlagWaterwalking                  MovementFlag = 0x10000000
	MovementFlagFallingSlow                   MovementFlag = 0x20000000
	MovementFlagHover                         MovementFlag = 0x40000000
	MovementFlagNoStrafe                      MovementFlag = 0x100000000
	MovementFlagNoJumping                     MovementFlag = 0x200000000
	MovementFlagUnknown1                      MovementFlag = 0x400000000
	MovementFlagFullSpeedTurning              MovementFlag = 0x800000000
	MovementFlagFullSpeedPitching             MovementFlag = 0x1000000000
	MovementFlagAlwaysAllowPitching           MovementFlag = 0x2000000000
	MovementFlagUnknown2                      MovementFlag = 0x4000000000
	MovementFlagUnknown3                      MovementFlag = 0x8000000000
	MovementFlagUnknown4                      MovementFlag = 0x10000000000
	MovementFlagUnknown5                      MovementFlag = 0x20000000000
	MovementFlagTransportInterpolatedMovement MovementFlag = 0x40000000000
	MovementFlagTransportInterpolatedTurning  MovementFlag = 0x80000000000
	MovementFlagTransportInterpolatedPitching MovementFlag = 0x100000000000
	MovementFlagUnknown6                      MovementFlag = 0x200000000000
	MovementFlagUnknown7                      MovementFlag = 0x400000000000
	MovementFlagUnknown8                      MovementFlag = 0x800000000000
)

// Required
// MovementUpdateLiving
type LivingFlags struct {
	Flags MovementFlag
}

// Required
// MovementUpdateLiving
type LivingCommonData1 struct {
	Timestamp        uint32
	PositionRotation realmd.Vector4
}

// Optional
// MovementFlagTransportInterpolatedMovement
// Nested in living block
type TransportPassengerInterpolatedData struct {
	// https://gtker.com/wow_messages/docs/transportinfo.html#wowm-representation-1
	// TODO: TransportInfo
	TransportTime uint32
}

// Optional
// MovementFlagOnTransport
// Mutually exclusive with MovementFlagTransportInterpolatedMovement
// Nested in living block
type TransportPassengerData struct {
	// TODO: TransportInfo
}

// Optional
// MovementFlagSwimming OR MovementFlagFlying OR MovementFlagAlwaysAllowPitching
// Nested in living block
type PitchData struct {
	Pitch float32
}

// Required
// MovementUpdateLiving
// Nested in living block
type LivingCommonData2 struct {
	// FallTime always is included even if MovementFlagFalling isn't set
	FallTime float32
}

// Optional
// MovementFlagFalling
// Nested in living block
type FallData struct {
	FallSpeed       float32
	CosAngle        float32 // TODO: name?
	SinAngle        float32 // TODO: name?
	HorizontalSpeed float32 // TODO: research
}

// Optional
// MovementFlagSplineElevation
// Nested in living block
type SplineElevationData struct {
	Elevation float32 // TODO: research
}

// Required
// MovementUpdateLiving
// Nested in living block
type LivingCommonData3 struct {
	WalkSpeed          float32
	RunSpeed           float32
	ReverseSpeed       float32
	SwimSpeed          float32
	SwimReverseSpeed   float32
	FlightSpeed        float32
	FlightReverseSpeed float32
	TurnRate           float32
	PitchRate          float32
}

type LivingDataCommon struct {
	LivingCommonData1
	LivingCommonData2
	LivingCommonData3
}

// Optional
// MovementFlagSplineElevation
type SplineData struct {
	// TODO
}

// MovementUpdatePosition
type PositionData struct {
	// TODO
}

// MovementUpdateHasPosition
type HasPositionData struct {
	// TODO
}

// MovementUpdateHighGuid
type HighGuidData struct {
	HighGuid uint32
}

// MovementUpdateLowGuid
type LowGuidData struct {
	LowGuid uint32
}

// MovementUpdateHasAttackingTarget
type AttackingTargetData struct {
	Guid realmd.PackedGuid
}

// MovementUpdateTransport
type TransportData struct {
	TransportProgressMs uint32 // milliseconds
}

// MovementUpdateVehicle
type VehicleData struct {
	Id          uint32
	Orientation float32
}

// MovementUpdateRotation
type RotationData struct {
	PackedLocalRotation uint64 // TODO: research
}

type movementBuffer struct {
	updateFlag MovementUpdateFlag

	// living block
	livingFlags                    MovementFlag
	living1                        *LivingCommonData1
	transportPassengerInterpolated *TransportPassengerInterpolatedData
	transportPassenger             *TransportPassengerData
	pitch                          *PitchData
	living2                        *LivingCommonData2
	fall                           *FallData
	splineElevation                *SplineElevationData
	living3                        *LivingCommonData3
	spline                         *SplineData

	positionData        *PositionData    // TODO: naming
	hasPositionData     *HasPositionData // TODO: naming
	highGuidData        *HighGuidData
	lowGuidData         *LowGuidData
	attackingTargetData *AttackingTargetData
	transportData       *TransportData
	vehicleData         *VehicleData
	rotationData        *RotationData
}

// MovementBuild builds the movement block.
// https://gtker.com/wow_messages/docs/movementblock.html#client-version-335
type MovementBuilder struct {
	buf           movementBuffer
	livingBuilder *LivingMovementBuilder
}

// Bytes returns the movement block as a little-endian byte array.
func (b *MovementBuilder) Bytes() []byte {
	bytesBuf := bytes.Buffer{}
	binary.Write(&bytesBuf, binary.LittleEndian, b.buf.updateFlag)

	if b.buf.updateFlag&MovementUpdateLiving > 0 {
		flags := b.buf.livingFlags

		binary.Write(&bytesBuf, binary.LittleEndian, flags)
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.living1)

		if b.buf.transportPassengerInterpolated != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.transportPassengerInterpolated)
		}
		if b.buf.transportPassenger != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.transportPassenger)
		}
		if b.buf.pitch != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.pitch)
		}

		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.living2)

		if b.buf.fall != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.fall)
		}
		if b.buf.splineElevation != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.splineElevation)
		}

		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.living3)

		if b.buf.spline != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, b.buf.spline)
		}
	} else if b.buf.updateFlag&MovementUpdatePosition > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.positionData)
	} else if b.buf.updateFlag&MovementUpdateHasPosition > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.hasPositionData)
	}

	if b.buf.updateFlag&MovementUpdateHighGuid > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.highGuidData)
	}
	if b.buf.updateFlag&MovementUpdateLowGuid > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.lowGuidData)
	}
	if b.buf.updateFlag&MovementUpdateHasAttackingTarget > 0 {
		bytesBuf.Write(b.buf.attackingTargetData.Guid)
	}
	if b.buf.updateFlag&MovementUpdateTransport > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.transportData)
	}
	if b.buf.updateFlag&MovementUpdateVehicle > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.vehicleData)
	}
	if b.buf.updateFlag&MovementUpdateRotation > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, b.buf.rotationData)
	}

	return bytesBuf.Bytes()
}

// Living returns a builder for building the block related to living units. It also sets the living flag.
func (b *MovementBuilder) Living() *LivingMovementBuilder {
	if b.livingBuilder == nil {
		b.livingBuilder = &LivingMovementBuilder{buf: &b.buf}
		b.buf.updateFlag |= MovementUpdateLiving
	}
	return b.livingBuilder
}

func (b *MovementBuilder) Position(data *PositionData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdatePosition
	} else {
		b.buf.updateFlag |= MovementUpdatePosition
	}
	b.buf.positionData = data
}

func (b *MovementBuilder) HasPosition(data *HasPositionData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateHasPosition
	} else {
		b.buf.updateFlag |= MovementUpdateHasPosition
	}
	b.buf.hasPositionData = data
}

func (b *MovementBuilder) HighGuid(data *HighGuidData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateHighGuid
	} else {
		b.buf.updateFlag |= MovementUpdateHighGuid
	}
	b.buf.highGuidData = data
}

func (b *MovementBuilder) LowGuid(data *LowGuidData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateLowGuid
	} else {
		b.buf.updateFlag |= MovementUpdateLowGuid
	}
	b.buf.lowGuidData = data
}

func (b *MovementBuilder) AttackingTarget(data *AttackingTargetData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateHasAttackingTarget
	} else {
		b.buf.updateFlag |= MovementUpdateHasAttackingTarget
	}
	b.buf.attackingTargetData = data
}

func (b *MovementBuilder) Transport(data *TransportData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateTransport
	} else {
		b.buf.updateFlag |= MovementUpdateTransport
	}
	b.buf.transportData = data
}

func (b *MovementBuilder) Vehicle(data *VehicleData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateVehicle
	} else {
		b.buf.updateFlag |= MovementUpdateVehicle
	}
	b.buf.vehicleData = data
}

func (b *MovementBuilder) Rotation(data *RotationData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateRotation
	} else {
		b.buf.updateFlag |= MovementUpdateRotation
	}
	b.buf.rotationData = data
}

type LivingMovementBuilder struct {
	buf *movementBuffer
}

// Common sets the fields always present in the living block. Panics if data is nil.
func (b *LivingMovementBuilder) Common(data *LivingDataCommon) error {
	if data == nil {
		panic("data cannot be nil")
	}

	b.buf.living1 = &data.LivingCommonData1
	b.buf.living2 = &data.LivingCommonData2
	b.buf.living3 = &data.LivingCommonData3

	return nil
}

func (b *LivingMovementBuilder) TransportPassengerMovement(data *TransportPassengerInterpolatedData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagTransportInterpolatedMovement
	} else {
		b.buf.livingFlags |= MovementFlagTransportInterpolatedMovement
	}
	b.buf.transportPassengerInterpolated = data
}

func (b *LivingMovementBuilder) TransportPassenger(data *TransportPassengerData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagOnTransport
	} else {
		b.buf.livingFlags |= MovementFlagOnTransport
	}
	b.buf.transportPassenger = data
}

func (b *LivingMovementBuilder) Swimming(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagSwimming
	} else {
		b.buf.livingFlags |= MovementFlagSwimming
	}
	b.buf.pitch = data
}

func (b *LivingMovementBuilder) Flying(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagFlying
	} else {
		b.buf.livingFlags |= MovementFlagFlying
	}
	b.buf.pitch = data
}

func (b *LivingMovementBuilder) ForcePitch(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagAlwaysAllowPitching
	} else {
		b.buf.livingFlags |= MovementFlagAlwaysAllowPitching
	}
	b.buf.pitch = data
}

func (b *LivingMovementBuilder) Falling(data *FallData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagFalling
	} else {
		b.buf.livingFlags |= MovementFlagFalling
	}
	b.buf.fall = data
}

func (b *LivingMovementBuilder) SplineElevation(data *SplineElevationData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagSplineElevation
	} else {
		b.buf.livingFlags |= MovementFlagSplineElevation
	}
	b.buf.splineElevation = data
}

func (b *LivingMovementBuilder) Spline(data *SplineData) {
	if data == nil {
		b.buf.livingFlags &= ^MovementFlagSplineEnabled
	} else {
		b.buf.livingFlags |= MovementFlagSplineEnabled
	}
	b.buf.spline = data
}
