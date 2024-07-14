package values

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
	MovementUpdateStationaryPosition MovementUpdateFlag = 0x40
	MovementUpdateVehicle            MovementUpdateFlag = 0x80
	MovementUpdatePosition           MovementUpdateFlag = 0x100
	MovementUpdateRotation           MovementUpdateFlag = 0x200
)

// Encoded as 48 bits
type LivingMovementFlag uint64

const (
	LivingMovementFlagNone                          LivingMovementFlag = 0x0
	LivingMovementFlagForward                       LivingMovementFlag = 0x1
	LivingMovementFlagBackward                      LivingMovementFlag = 0x2
	LivingMovementFlagStrafeLeft                    LivingMovementFlag = 0x4
	LivingMovementFlagStrafeRight                   LivingMovementFlag = 0x8
	LivingMovementFlagLeft                          LivingMovementFlag = 0x10
	LivingMovementFlagRight                         LivingMovementFlag = 0x20
	LivingMovementFlagPitchUp                       LivingMovementFlag = 0x40
	LivingMovementFlagPitchDown                     LivingMovementFlag = 0x80
	LivingMovementFlagWalking                       LivingMovementFlag = 0x100
	LivingMovementFlagOnTransport                   LivingMovementFlag = 0x200
	LivingMovementFlagDisableGravity                LivingMovementFlag = 0x400
	LivingMovementFlagRoot                          LivingMovementFlag = 0x800
	LivingMovementFlagFalling                       LivingMovementFlag = 0x1000
	LivingMovementFlagFallingFar                    LivingMovementFlag = 0x2000
	LivingMovementFlagPendingStop                   LivingMovementFlag = 0x4000
	LivingMovementFlagPendingStrafeStop             LivingMovementFlag = 0x8000
	LivingMovementFlagPendingForward                LivingMovementFlag = 0x10000
	LivingMovementFlagPendingBackward               LivingMovementFlag = 0x20000
	LivingMovementFlagPendingStrafeLeft             LivingMovementFlag = 0x40000
	LivingMovementFlagPendingStrafeRight            LivingMovementFlag = 0x80000
	LivingMovementFlagPendingRoot                   LivingMovementFlag = 0x100000
	LivingMovementFlagSwimming                      LivingMovementFlag = 0x200000
	LivingMovementFlagAscending                     LivingMovementFlag = 0x400000
	LivingMovementFlagDescending                    LivingMovementFlag = 0x800000
	LivingMovementFlagCanFly                        LivingMovementFlag = 0x1000000
	LivingMovementFlagFlying                        LivingMovementFlag = 0x2000000
	LivingMovementFlagSplineElevation               LivingMovementFlag = 0x4000000
	LivingMovementFlagSplineEnabled                 LivingMovementFlag = 0x8000000
	LivingMovementFlagWaterwalking                  LivingMovementFlag = 0x10000000
	LivingMovementFlagFallingSlow                   LivingMovementFlag = 0x20000000
	LivingMovementFlagHover                         LivingMovementFlag = 0x40000000
	LivingMovementFlagNoStrafe                      LivingMovementFlag = 0x100000000
	LivingMovementFlagNoJumping                     LivingMovementFlag = 0x200000000
	LivingMovementFlagUnknown1                      LivingMovementFlag = 0x400000000
	LivingMovementFlagFullSpeedTurning              LivingMovementFlag = 0x800000000
	LivingMovementFlagFullSpeedPitching             LivingMovementFlag = 0x1000000000
	LivingMovementFlagAlwaysAllowPitching           LivingMovementFlag = 0x2000000000
	LivingMovementFlagUnknown2                      LivingMovementFlag = 0x4000000000
	LivingMovementFlagUnknown3                      LivingMovementFlag = 0x8000000000
	LivingMovementFlagUnknown4                      LivingMovementFlag = 0x10000000000
	LivingMovementFlagUnknown5                      LivingMovementFlag = 0x20000000000
	LivingMovementFlagTransportInterpolatedMovement LivingMovementFlag = 0x40000000000
	LivingMovementFlagTransportInterpolatedTurning  LivingMovementFlag = 0x80000000000
	LivingMovementFlagTransportInterpolatedPitching LivingMovementFlag = 0x100000000000
	LivingMovementFlagUnknown6                      LivingMovementFlag = 0x200000000000
	LivingMovementFlagTransitionBetweenSwimAndFly   LivingMovementFlag = 0x400000000000
	LivingMovementFlagUnknown8                      LivingMovementFlag = 0x800000000000
)

// LivingData contains data that is always present when the living flag is set.
type LivingData struct {
	// First section
	Timestamp        uint32
	PositionRotation realmd.Vector4

	// Second section
	// FallTime always is included even if LivingMovementFlagFalling isn't set
	FallTime float32

	// Third section
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

// MovementUpdateLiving (optional)
type TransportPassengerInterpolatedData struct {
	// https://gtker.com/wow_messages/docs/transportinfo.html#wowm-representation-1
	// TODO: TransportInfo
	TransportTime uint32
}

// MovementUpdateLiving (optional)
type TransportPassengerData struct {
	// TODO: TransportInfo
}

// MovementUpdateLiving (optional)
type PitchData struct {
	Pitch float32
}

// MovementUpdateLiving (optional)
type FallData struct {
	FallSpeed       float32
	CosAngle        float32 // TODO: name?
	SinAngle        float32 // TODO: name?
	HorizontalSpeed float32 // TODO: research
}

// MovementUpdateLiving (optional)
type SplineElevationData struct {
	Elevation float32 // TODO: research
}

// MovementUpdateLiving (optional)
type SplineData struct {
	// TODO
}

// MovementUpdatePosition
type PositionData struct {
	// TODO
}

// MovementUpdateStationaryPosition
type StationaryPositionData struct {
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
	livingFlags                    LivingMovementFlag
	livingData                     *LivingData
	transportPassengerInterpolated *TransportPassengerInterpolatedData
	transportPassenger             *TransportPassengerData
	pitch                          *PitchData
	fall                           *FallData
	splineElevation                *SplineElevationData
	spline                         *SplineData

	positionData        *PositionData           // TODO: naming
	hasPositionData     *StationaryPositionData // TODO: naming
	highGuidData        *HighGuidData
	lowGuidData         *LowGuidData
	attackingTargetData *AttackingTargetData
	transportData       *TransportData
	vehicleData         *VehicleData
	rotationData        *RotationData
}

// MovementValues provides an interface for setting values related to movement and position.
// Values can be added in any order.
// https://gtker.com/wow_messages/docs/movementblock.html#client-version-335
type MovementValues struct {
	buf           movementBuffer
	livingBuilder *LivingMovementBuilder
}

// Bytes returns the movement block as a little-endian byte array.
func (m *MovementValues) Bytes() []byte {
	bytesBuf := bytes.Buffer{}
	binary.Write(&bytesBuf, binary.LittleEndian, m.buf.updateFlag)

	if m.buf.updateFlag&MovementUpdateLiving > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingFlags)

		// The flag is 6 bytes but stored as 8, discard the last 2 bytes.
		bytesBuf.Truncate(bytesBuf.Len() - 2)

		// First section of common data
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.Timestamp)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.PositionRotation)

		if m.buf.transportPassengerInterpolated != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.transportPassengerInterpolated)
		}
		if m.buf.transportPassenger != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.transportPassenger)
		}
		if m.buf.pitch != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.pitch)
		}

		// Second section of common data
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.FallTime)

		if m.buf.fall != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.fall)
		}
		if m.buf.splineElevation != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.splineElevation)
		}

		// Third section of common data
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.WalkSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.RunSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.ReverseSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.SwimSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.SwimReverseSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.FlightSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.FlightReverseSpeed)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.TurnRate)
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.livingData.PitchRate)

		if m.buf.spline != nil {
			binary.Write(&bytesBuf, binary.LittleEndian, m.buf.spline)
		}
	} else if m.buf.updateFlag&MovementUpdatePosition > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.positionData)
	} else if m.buf.updateFlag&MovementUpdateStationaryPosition > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.hasPositionData)
	}

	if m.buf.updateFlag&MovementUpdateHighGuid > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.highGuidData)
	}
	if m.buf.updateFlag&MovementUpdateLowGuid > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.lowGuidData)
	}
	if m.buf.updateFlag&MovementUpdateHasAttackingTarget > 0 {
		bytesBuf.Write(m.buf.attackingTargetData.Guid)
	}
	if m.buf.updateFlag&MovementUpdateTransport > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.transportData)
	}
	if m.buf.updateFlag&MovementUpdateVehicle > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.vehicleData)
	}
	if m.buf.updateFlag&MovementUpdateRotation > 0 {
		binary.Write(&bytesBuf, binary.LittleEndian, m.buf.rotationData)
	}

	return bytesBuf.Bytes()
}

// Living returns a builder for building the block related to living units. It also sets the living flag.
func (m *MovementValues) Living() *LivingMovementBuilder {
	if m.livingBuilder == nil {
		m.livingBuilder = &LivingMovementBuilder{buf: &m.buf}
		m.buf.updateFlag |= MovementUpdateLiving
	}

	return m.livingBuilder
}

func (m *MovementValues) Self() {
	m.buf.updateFlag |= MovementUpdateSelf
}

func (m *MovementValues) Position(data *PositionData) {
	if data == nil {
		m.buf.updateFlag &= ^MovementUpdatePosition
	} else {
		m.buf.updateFlag |= MovementUpdatePosition
	}

	m.buf.positionData = data
}

func (m *MovementValues) HasPosition(data *StationaryPositionData) {
	if data == nil {
		m.buf.updateFlag &= ^MovementUpdateStationaryPosition
	} else {
		m.buf.updateFlag |= MovementUpdateStationaryPosition
	}

	m.buf.hasPositionData = data
}

func (b *MovementValues) HighGuid(data *HighGuidData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateHighGuid
	} else {
		b.buf.updateFlag |= MovementUpdateHighGuid
	}

	b.buf.highGuidData = data
}

func (b *MovementValues) LowGuid(data *LowGuidData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateLowGuid
	} else {
		b.buf.updateFlag |= MovementUpdateLowGuid
	}

	b.buf.lowGuidData = data
}

func (b *MovementValues) AttackingTarget(data *AttackingTargetData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateHasAttackingTarget
	} else {
		b.buf.updateFlag |= MovementUpdateHasAttackingTarget
	}

	b.buf.attackingTargetData = data
}

func (b *MovementValues) Transport(data *TransportData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateTransport
	} else {
		b.buf.updateFlag |= MovementUpdateTransport
	}

	b.buf.transportData = data
}

func (b *MovementValues) Vehicle(data *VehicleData) {
	if data == nil {
		b.buf.updateFlag &= ^MovementUpdateVehicle
	} else {
		b.buf.updateFlag |= MovementUpdateVehicle
	}

	b.buf.vehicleData = data
}

func (b *MovementValues) Rotation(data *RotationData) {
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

// Data sets the fields always present in the living block. Panics if data is nil.
func (b *LivingMovementBuilder) Data(data *LivingData) error {
	if data == nil {
		panic("data cannot be nil")
	}

	b.buf.livingData = data

	return nil
}

func (b *LivingMovementBuilder) TransportPassengerMovement(data *TransportPassengerInterpolatedData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagTransportInterpolatedMovement
	} else {
		b.buf.livingFlags |= LivingMovementFlagTransportInterpolatedMovement
	}

	b.buf.transportPassengerInterpolated = data
}

func (b *LivingMovementBuilder) TransportPassenger(data *TransportPassengerData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagOnTransport
	} else {
		b.buf.livingFlags |= LivingMovementFlagOnTransport
	}

	b.buf.transportPassenger = data
}

func (b *LivingMovementBuilder) Swimming(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagSwimming
	} else {
		b.buf.livingFlags |= LivingMovementFlagSwimming
	}

	b.buf.pitch = data
}

func (b *LivingMovementBuilder) Flying(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagFlying
	} else {
		b.buf.livingFlags |= LivingMovementFlagFlying
	}

	b.buf.pitch = data
}

func (b *LivingMovementBuilder) ForcePitch(data *PitchData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagAlwaysAllowPitching
	} else {
		b.buf.livingFlags |= LivingMovementFlagAlwaysAllowPitching
	}

	b.buf.pitch = data
}

func (b *LivingMovementBuilder) Falling(data *FallData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagFalling
	} else {
		b.buf.livingFlags |= LivingMovementFlagFalling
	}

	b.buf.fall = data
}

func (b *LivingMovementBuilder) SplineElevation(data *SplineElevationData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagSplineElevation
	} else {
		b.buf.livingFlags |= LivingMovementFlagSplineElevation
	}

	b.buf.splineElevation = data
}

func (b *LivingMovementBuilder) Spline(data *SplineData) {
	if data == nil {
		b.buf.livingFlags &= ^LivingMovementFlagSplineEnabled
	} else {
		b.buf.livingFlags |= LivingMovementFlagSplineEnabled
	}

	b.buf.spline = data
}
