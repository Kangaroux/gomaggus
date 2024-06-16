package objupdate

// This is encoded as 48 bits
type MovementFlag uint64

const (
	MovementFlagNone                 MovementFlag = 0x0
	MovementFlagForward              MovementFlag = 0x1
	MovementFlagBackward             MovementFlag = 0x2
	MovementFlagStrafeLeft           MovementFlag = 0x4
	MovementFlagStrafeRight          MovementFlag = 0x8
	MovementFlagLeft                 MovementFlag = 0x10
	MovementFlagRight                MovementFlag = 0x20
	MovementFlagPitchUp              MovementFlag = 0x40
	MovementFlagPitchDown            MovementFlag = 0x80
	MovementFlagWalking              MovementFlag = 0x100
	MovementFlagOnTransport          MovementFlag = 0x200
	MovementFlagDisableGravity       MovementFlag = 0x400
	MovementFlagRoot                 MovementFlag = 0x800
	MovementFlagFalling              MovementFlag = 0x1000
	MovementFlagFallingFar           MovementFlag = 0x2000
	MovementFlagPendingStop          MovementFlag = 0x4000
	MovementFlagPendingStrafeStop    MovementFlag = 0x8000
	MovementFlagPendingForward       MovementFlag = 0x10000
	MovementFlagPendingBackward      MovementFlag = 0x20000
	MovementFlagPendingStrafeLeft    MovementFlag = 0x40000
	MovementFlagPendingStrafeRight   MovementFlag = 0x80000
	MovementFlagPendingRoot          MovementFlag = 0x100000
	MovementFlagSwimming             MovementFlag = 0x200000
	MovementFlagAscending            MovementFlag = 0x400000
	MovementFlagDescending           MovementFlag = 0x800000
	MovementFlagCanFly               MovementFlag = 0x1000000
	MovementFlagFlying               MovementFlag = 0x2000000
	MovementFlagSplineElevation      MovementFlag = 0x4000000
	MovementFlagSplineEnabled        MovementFlag = 0x8000000
	MovementFlagWaterwalking         MovementFlag = 0x10000000
	MovementFlagFallingSlow          MovementFlag = 0x20000000
	MovementFlagHover                MovementFlag = 0x40000000
	MovementFlagNoStrafe             MovementFlag = 0x100000000
	MovementFlagNoJumping            MovementFlag = 0x200000000
	MovementFlagUnknown1             MovementFlag = 0x400000000
	MovementFlagFullSpeedTurning     MovementFlag = 0x800000000
	MovementFlagFullSpeedPitching    MovementFlag = 0x1000000000
	MovementFlagAlwaysAllowPitching  MovementFlag = 0x2000000000
	MovementFlagUnknown2             MovementFlag = 0x4000000000
	MovementFlagUnknown3             MovementFlag = 0x8000000000
	MovementFlagUnknown4             MovementFlag = 0x10000000000
	MovementFlagUnknown5             MovementFlag = 0x20000000000
	MovementFlagInterpolatedMovement MovementFlag = 0x40000000000
	MovementFlagInterpolatedTurning  MovementFlag = 0x80000000000
	MovementFlagInterpolatedPitching MovementFlag = 0x100000000000
	MovementFlagUnknown6             MovementFlag = 0x200000000000
	MovementFlagUnknown7             MovementFlag = 0x400000000000
	MovementFlagUnknown8             MovementFlag = 0x800000000000
)
