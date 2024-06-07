package worldd

// Opcodes sent by the server
const (
	OP_SRV_AUTH_CHALLENGE          uint16 = 0x1EC
	OP_SRV_AUTH_RESPONSE           uint16 = 0x1EE
	OP_SRV_PONG                    uint16 = 0x1DD
	OP_SRV_ACCOUNT_DATA_TIMES      uint16 = 0x209
	OP_SRV_CHAR_ENUM               uint16 = 0x3B
	OP_SRV_REALM_SPLIT             uint16 = 0x38B
	OP_SRV_CHAR_CREATE             uint16 = 0x3A
	OP_SRV_CHAR_DELETE             uint16 = 0x3C
	OP_SRV_CHAR_LOGIN_FAILED       uint16 = 0x41
	OP_SRV_CHAR_LOGIN_VERIFY_WORLD uint16 = 0x236
	OP_SRV_UPDATE_OBJECT           uint16 = 0xA9
	OP_SRV_TUTORIAL_FLAGS          uint16 = 0xFD
)

// Opcodes sent by the client
const (
	OP_CL_AUTH_SESSION uint32 = 0x1ED
	OP_CL_REALM_SPLIT  uint32 = 0x38C
	OP_CL_PING         uint32 = 0x1DC
	OP_CL_CHAR_ENUM    uint32 = 0x37
	OP_CL_CHAR_CREATE  uint32 = 0x36
	OP_CL_CHAR_DELETE  uint32 = 0x38
	OP_CL_PLAYER_LOGIN uint32 = 0x3D

	// Client sent after receiving our OP_SRV_AUTH_RESPONSE. The packet is empty besides the header.
	// Immediately after the client sends OP_CL_CHAR_ENUM.
	OP_CL_READY_FOR_ACCOUNT_DATA_TIMES uint32 = 0x4FF
)

type ResponseCode = byte

const (
	RespCodeResponseSuccess                                ResponseCode = 0x00
	RespCodeResponseFailure                                ResponseCode = 0x01
	RespCodeResponseCancelled                              ResponseCode = 0x02
	RespCodeResponseDisconnected                           ResponseCode = 0x03
	RespCodeResponseFailedToConnect                        ResponseCode = 0x04
	RespCodeResponseConnected                              ResponseCode = 0x05
	RespCodeResponseVersionMismatch                        ResponseCode = 0x06
	RespCodeCStatusConnecting                              ResponseCode = 0x07
	RespCodeCStatusNegotiatingSecurity                     ResponseCode = 0x08
	RespCodeCStatusNegotiationComplete                     ResponseCode = 0x09
	RespCodeCStatusNegotiationFailed                       ResponseCode = 0x0A
	RespCodeCStatusAuthenticating                          ResponseCode = 0x0B
	RespCodeAuthOk                                         ResponseCode = 0x0C
	RespCodeAuthFailed                                     ResponseCode = 0x0D
	RespCodeAuthReject                                     ResponseCode = 0x0E
	RespCodeAuthBadServerProof                             ResponseCode = 0x0F
	RespCodeAuthUnavailable                                ResponseCode = 0x10
	RespCodeAuthSystemError                                ResponseCode = 0x11
	RespCodeAuthBillingError                               ResponseCode = 0x12
	RespCodeAuthBillingExpired                             ResponseCode = 0x13
	RespCodeAuthVersionMismatch                            ResponseCode = 0x14
	RespCodeAuthUnknownAccount                             ResponseCode = 0x15
	RespCodeAuthIncorrectPassword                          ResponseCode = 0x16
	RespCodeAuthSessionExpired                             ResponseCode = 0x17
	RespCodeAuthServerShuttingDown                         ResponseCode = 0x18
	RespCodeAuthAlreadyLoggingIn                           ResponseCode = 0x19
	RespCodeAuthLoginServerNotFound                        ResponseCode = 0x1A
	RespCodeAuthWaitQueue                                  ResponseCode = 0x1B
	RespCodeAuthBanned                                     ResponseCode = 0x1C
	RespCodeAuthAlreadyOnline                              ResponseCode = 0x1D
	RespCodeAuthNoTime                                     ResponseCode = 0x1E
	RespCodeAuthDbBusy                                     ResponseCode = 0x1F
	RespCodeAuthSuspended                                  ResponseCode = 0x20
	RespCodeAuthParentalControl                            ResponseCode = 0x21
	RespCodeAuthLockedEnforced                             ResponseCode = 0x22
	RespCodeRealmListInProgress                            ResponseCode = 0x23
	RespCodeRealmListSuccess                               ResponseCode = 0x24
	RespCodeRealmListFailed                                ResponseCode = 0x25
	RespCodeRealmListInvalid                               ResponseCode = 0x26
	RespCodeRealmListRealmNotFound                         ResponseCode = 0x27
	RespCodeAccountCreateInProgress                        ResponseCode = 0x28
	RespCodeAccountCreateSuccess                           ResponseCode = 0x29
	RespCodeAccountCreateFailed                            ResponseCode = 0x2A
	RespCodeCharListRetrieving                             ResponseCode = 0x2B
	RespCodeCharListRetrieved                              ResponseCode = 0x2C
	RespCodeCharListFailed                                 ResponseCode = 0x2D
	RespCodeCharCreateInProgress                           ResponseCode = 0x2E
	RespCodeCharCreateSuccess                              ResponseCode = 0x2F
	RespCodeCharCreateError                                ResponseCode = 0x30
	RespCodeCharCreateFailed                               ResponseCode = 0x31
	RespCodeCharCreateNameInUse                            ResponseCode = 0x32
	RespCodeCharCreateDisabled                             ResponseCode = 0x33
	RespCodeCharCreatePvpTeamsViolation                    ResponseCode = 0x34
	RespCodeCharCreateServerLimit                          ResponseCode = 0x35
	RespCodeCharCreateAccountLimit                         ResponseCode = 0x36
	RespCodeCharCreateServerQueue                          ResponseCode = 0x37
	RespCodeCharCreateOnlyExisting                         ResponseCode = 0x38
	RespCodeCharCreateExpansion                            ResponseCode = 0x39
	RespCodeCharCreateExpansionClass                       ResponseCode = 0x3A
	RespCodeCharCreateLevelRequirement                     ResponseCode = 0x3B
	RespCodeCharCreateUniqueClassLimit                     ResponseCode = 0x3C
	RespCodeCharCreateCharacterInGuild                     ResponseCode = 0x3D
	RespCodeCharCreateRestrictedRaceclass                  ResponseCode = 0x3E
	RespCodeCharCreateCharacterChooseRace                  ResponseCode = 0x3F
	RespCodeCharCreateCharacterArenaLeader                 ResponseCode = 0x40
	RespCodeCharCreateCharacterDeleteMail                  ResponseCode = 0x41
	RespCodeCharCreateCharacterSwapFaction                 ResponseCode = 0x42
	RespCodeCharCreateCharacterRaceOnly                    ResponseCode = 0x43
	RespCodeCharCreateCharacterGoldLimit                   ResponseCode = 0x44
	RespCodeCharCreateForceLogin                           ResponseCode = 0x45
	RespCodeCharDeleteInProgress                           ResponseCode = 0x46
	RespCodeCharDeleteSuccess                              ResponseCode = 0x47
	RespCodeCharDeleteFailed                               ResponseCode = 0x48
	RespCodeCharDeleteFailedLockedForTransfer              ResponseCode = 0x49
	RespCodeCharDeleteFailedGuildLeader                    ResponseCode = 0x4A
	RespCodeCharDeleteFailedArenaCaptain                   ResponseCode = 0x4B
	RespCodeCharLoginInProgress                            ResponseCode = 0x4C
	RespCodeCharLoginSuccess                               ResponseCode = 0x4D
	RespCodeCharLoginNoWorld                               ResponseCode = 0x4E
	RespCodeCharLoginDuplicateCharacter                    ResponseCode = 0x4F
	RespCodeCharLoginNoInstances                           ResponseCode = 0x50
	RespCodeCharLoginFailed                                ResponseCode = 0x51
	RespCodeCharLoginDisabled                              ResponseCode = 0x52
	RespCodeCharLoginNoCharacter                           ResponseCode = 0x53
	RespCodeCharLoginLockedForTransfer                     ResponseCode = 0x54
	RespCodeCharLoginLockedByBilling                       ResponseCode = 0x55
	RespCodeCharLoginLockedByMobileAh                      ResponseCode = 0x56
	RespCodeCharNameSuccess                                ResponseCode = 0x57
	RespCodeCharNameFailure                                ResponseCode = 0x58
	RespCodeCharNameNoName                                 ResponseCode = 0x59
	RespCodeCharNameTooShort                               ResponseCode = 0x5A
	RespCodeCharNameTooLong                                ResponseCode = 0x5B
	RespCodeCharNameInvalidCharacter                       ResponseCode = 0x5C
	RespCodeCharNameMixedLanguages                         ResponseCode = 0x5D
	RespCodeCharNameProfane                                ResponseCode = 0x5E
	RespCodeCharNameReserved                               ResponseCode = 0x5F
	RespCodeCharNameInvalidApostrophe                      ResponseCode = 0x60
	RespCodeCharNameMultipleApostrophes                    ResponseCode = 0x61
	RespCodeCharNameThreeConsecutive                       ResponseCode = 0x62
	RespCodeCharNameInvalidSpace                           ResponseCode = 0x63
	RespCodeCharNameConsecutiveSpaces                      ResponseCode = 0x64
	RespCodeCharNameRussianConsecutiveSilentCharacters     ResponseCode = 0x65
	RespCodeCharNameRussianSilentCharacterAtBeginningOrEnd ResponseCode = 0x66
	RespCodeCharNameDeclensionDoesntMatchBaseName          ResponseCode = 0x67
)

type Expansion = byte

const (
	ExpansionVanilla Expansion = 0x0
	ExpansionTbc     Expansion = 0x1
	ExpansionWrath   Expansion = 0x2
)

type RealmSplitState = uint32

const (
	SplitNormal    = 0
	SplitConfirmed = 1
	SplitPending   = 2
)

type PowerType = byte

const (
	PowerTypeMana      PowerType = 0
	PowerTypeRage      PowerType = 1
	PowerTypeFocus     PowerType = 2
	PowerTypeEnergy    PowerType = 3
	PowerTypeHappiness PowerType = 4
)

type Race = byte

const (
	RaceHuman             Race = 1
	RaceOrc               Race = 2
	RaceDwarf             Race = 3
	RaceNightElf          Race = 4
	RaceUndead            Race = 5
	RaceTauren            Race = 6
	RaceGnome             Race = 7
	RaceTroll             Race = 8
	RaceGoblin            Race = 9
	RaceBlood_elf         Race = 10
	RaceDraenei           Race = 11
	RaceFelOrc            Race = 12
	RaceNaga              Race = 13
	RaceBroken            Race = 14
	RaceSkeleton          Race = 15
	RaceVrykul            Race = 16
	RaceTuskarr           Race = 17
	RaceForestTroll       Race = 18
	RaceTaunka            Race = 19
	RaceNorthrendSkeleton Race = 20
	RaceIceTroll          Race = 21
)

type Class = byte

const (
	ClassWarrior Class = 1
	ClassPaladin Class = 2
	ClassHunter  Class = 3
	ClassRogue   Class = 4
	ClassPriest  Class = 5
	// ClassDeathKnight Class = 6
	ClassShaman  Class = 7
	ClassMage    Class = 8
	ClassWarlock Class = 9
	ClassDruid   Class = 11
)

type Gender = byte

const (
	GenderMale   Gender = 0
	GenderFemale Gender = 1
	GenderNone   Gender = 2 // used by pets?
)

type ObjectType = byte

const (
	ObjectTypeObject        ObjectType = 0
	ObjectTypeItem          ObjectType = 1
	ObjectTypeContainer     ObjectType = 2
	ObjectTypeUnit          ObjectType = 3
	ObjectTypePlayer        ObjectType = 4
	ObjectTypeGameObject    ObjectType = 5
	ObjectTypeDynamicObject ObjectType = 6
	ObjectTypeCorpse        ObjectType = 7
)

type UpdateType = byte

const (
	UpdateTypePartial           UpdateType = 0
	UpdateTypeMovement          UpdateType = 1
	UpdateTypeCreateObject      UpdateType = 2
	UpdateTypeCreateObject2     UpdateType = 3
	UpdateTypeOutOfRangeObjects UpdateType = 4
	UpdateTypeNearObjects       UpdateType = 5
)

// https://gtker.com/wow_messages/docs/updateflag.html#client-version-335
type UpdateFlag = uint16

const (
	UpdateFlagNone               UpdateFlag = 0x0000
	UpdateFlagSelf               UpdateFlag = 0x0001
	UpdateFlagTransport          UpdateFlag = 0x0002
	UpdateFlagHasAttackingTarget UpdateFlag = 0x0004
	UpdateFlagLowGuid            UpdateFlag = 0x0008
	UpdateFlagHighGuid           UpdateFlag = 0x0010
	UpdateFlagLiving             UpdateFlag = 0x0020
	UpdateFlagHasPosition        UpdateFlag = 0x0040
	UpdateFlagVehicle            UpdateFlag = 0x0080
	UpdateFlagPosition           UpdateFlag = 0x0100
	UpdateFlagRotation           UpdateFlag = 0x0200
)

// This is encoded as 48 bits
type MovementFlag = uint64

const (
	MovementFlagNone                 MovementFlag = 0x00000000
	MovementFlagForward              MovementFlag = 0x00000001
	MovementFlagBackward             MovementFlag = 0x00000002
	MovementFlagStrafeLeft           MovementFlag = 0x00000004
	MovementFlagStrafeRight          MovementFlag = 0x00000008
	MovementFlagLeft                 MovementFlag = 0x00000010
	MovementFlagRight                MovementFlag = 0x00000020
	MovementFlagPitchUp              MovementFlag = 0x00000040
	MovementFlagPitchDown            MovementFlag = 0x00000080
	MovementFlagWalking              MovementFlag = 0x00000100
	MovementFlagOnTransport          MovementFlag = 0x00000200
	MovementFlagDisableGravity       MovementFlag = 0x00000400
	MovementFlagRoot                 MovementFlag = 0x00000800
	MovementFlagFalling              MovementFlag = 0x00001000
	MovementFlagFallingFar           MovementFlag = 0x00002000
	MovementFlagPendingStop          MovementFlag = 0x00004000
	MovementFlagPendingStrafeStop    MovementFlag = 0x00008000
	MovementFlagPendingForward       MovementFlag = 0x00010000
	MovementFlagPendingBackward      MovementFlag = 0x00020000
	MovementFlagPendingStrafeLeft    MovementFlag = 0x00040000
	MovementFlagPendingStrafeRight   MovementFlag = 0x00080000
	MovementFlagPendingRoot          MovementFlag = 0x00100000
	MovementFlagSwimming             MovementFlag = 0x00200000
	MovementFlagAscending            MovementFlag = 0x00400000
	MovementFlagDescending           MovementFlag = 0x00800000
	MovementFlagCanFly               MovementFlag = 0x01000000
	MovementFlagFlying               MovementFlag = 0x02000000
	MovementFlagSplineElevation      MovementFlag = 0x04000000
	MovementFlagSplineEnabled        MovementFlag = 0x08000000
	MovementFlagWaterwalking         MovementFlag = 0x10000000
	MovementFlagFallingSlow          MovementFlag = 0x20000000
	MovementFlagHover                MovementFlag = 0x40000000
	MovementFlagNoStrafe             MovementFlag = 0x0000000100000000
	MovementFlagNoJumping            MovementFlag = 0x0000000200000000
	MovementFlagUnknown1             MovementFlag = 0x0000000400000000
	MovementFlagFullSpeedTurning     MovementFlag = 0x0000000800000000
	MovementFlagFullSpeedPitching    MovementFlag = 0x0000001000000000
	MovementFlagAlwaysAllowPitching  MovementFlag = 0x0000002000000000
	MovementFlagUnknown2             MovementFlag = 0x0000004000000000
	MovementFlagUnknown3             MovementFlag = 0x0000008000000000
	MovementFlagUnknown4             MovementFlag = 0x0000010000000000
	MovementFlagUnknown5             MovementFlag = 0x0000020000000000
	MovementFlagInterpolatedMovement MovementFlag = 0x0000040000000000
	MovementFlagInterpolatedTurning  MovementFlag = 0x0000080000000000
	MovementFlagInterpolatedPitching MovementFlag = 0x0000100000000000
	MovementFlagUnknown6             MovementFlag = 0x0000200000000000
	MovementFlagUnknown7             MovementFlag = 0x0000400000000000
	MovementFlagUnknown8             MovementFlag = 0x0000800000000000
)

// https://gtker.com/wow_messages/types/update-mask.html#version-335
type FieldMask struct {
	// Size is the number of uint32 blocks used for the field data.
	Size uint8
	// Offset is the bit number that is set to indicate this field is included.
	Offset uint16
}

var (
	FieldMaskObjectGuid    = FieldMask{Size: 2, Offset: 0x0000}
	FieldMaskObjectType    = FieldMask{Size: 1, Offset: 0x0002}
	FieldMaskObjectEntry   = FieldMask{Size: 1, Offset: 0x0003}
	FieldMaskObjectScaleX  = FieldMask{Size: 1, Offset: 0x0004}
	FieldMaskObjectPadding = FieldMask{Size: 1, Offset: 0x0005}

	FieldMaskItemOwner              = FieldMask{Size: 2, Offset: 0x0006}
	FieldMaskItemContained          = FieldMask{Size: 2, Offset: 0x0008}
	FieldMaskItemCreator            = FieldMask{Size: 2, Offset: 0x000A}
	FieldMaskItemGiftcreator        = FieldMask{Size: 2, Offset: 0x000C}
	FieldMaskItemStackCount         = FieldMask{Size: 1, Offset: 0x000E}
	FieldMaskItemDuration           = FieldMask{Size: 1, Offset: 0x000F}
	FieldMaskItemSpellCharges       = FieldMask{Size: 5, Offset: 0x0010}
	FieldMaskItemFlags              = FieldMask{Size: 1, Offset: 0x0015}
	FieldMaskItemEnchantment1_1     = FieldMask{Size: 2, Offset: 0x0016}
	FieldMaskItemEnchantment1_3     = FieldMask{Size: 1, Offset: 0x0018}
	FieldMaskItemEnchantment2_1     = FieldMask{Size: 2, Offset: 0x0019}
	FieldMaskItemEnchantment2_3     = FieldMask{Size: 1, Offset: 0x001B}
	FieldMaskItemEnchantment3_1     = FieldMask{Size: 2, Offset: 0x001C}
	FieldMaskItemEnchantment3_3     = FieldMask{Size: 1, Offset: 0x001E}
	FieldMaskItemEnchantment4_1     = FieldMask{Size: 2, Offset: 0x001F}
	FieldMaskItemEnchantment4_3     = FieldMask{Size: 1, Offset: 0x0021}
	FieldMaskItemEnchantment5_1     = FieldMask{Size: 2, Offset: 0x0022}
	FieldMaskItemEnchantment5_3     = FieldMask{Size: 1, Offset: 0x0024}
	FieldMaskItemEnchantment6_1     = FieldMask{Size: 2, Offset: 0x0025}
	FieldMaskItemEnchantment6_3     = FieldMask{Size: 1, Offset: 0x0027}
	FieldMaskItemEnchantment7_1     = FieldMask{Size: 2, Offset: 0x0028}
	FieldMaskItemEnchantment7_3     = FieldMask{Size: 1, Offset: 0x002A}
	FieldMaskItemEnchantment8_1     = FieldMask{Size: 2, Offset: 0x002B}
	FieldMaskItemEnchantment8_3     = FieldMask{Size: 1, Offset: 0x002D}
	FieldMaskItemEnchantment9_1     = FieldMask{Size: 2, Offset: 0x002E}
	FieldMaskItemEnchantment9_3     = FieldMask{Size: 1, Offset: 0x0030}
	FieldMaskItemEnchantment10_1    = FieldMask{Size: 2, Offset: 0x0031}
	FieldMaskItemEnchantment10_3    = FieldMask{Size: 1, Offset: 0x0033}
	FieldMaskItemEnchantment11_1    = FieldMask{Size: 2, Offset: 0x0034}
	FieldMaskItemEnchantment11_3    = FieldMask{Size: 1, Offset: 0x0036}
	FieldMaskItemEnchantment12_1    = FieldMask{Size: 2, Offset: 0x0037}
	FieldMaskItemEnchantment12_3    = FieldMask{Size: 1, Offset: 0x0039}
	FieldMaskItemPropertySeed       = FieldMask{Size: 1, Offset: 0x003A}
	FieldMaskItemRandomPropertiesId = FieldMask{Size: 1, Offset: 0x003B}
	FieldMaskItemDurability         = FieldMask{Size: 1, Offset: 0x003C}
	FieldMaskItemMaxdurability      = FieldMask{Size: 1, Offset: 0x003D}
	FieldMaskItemCreatePlayedTime   = FieldMask{Size: 1, Offset: 0x003E}

	FieldMaskContainerNumSlots = FieldMask{Size: 1, Offset: 0x0040}
	FieldMaskContainerSlot1    = FieldMask{Size: 7, Offset: 0x0042}

	FieldMaskUnitCharm                             = FieldMask{Size: 2, Offset: 0x0006}
	FieldMaskUnitSummon                            = FieldMask{Size: 2, Offset: 0x0008}
	FieldMaskUnitCritter                           = FieldMask{Size: 2, Offset: 0x000A}
	FieldMaskUnitCharmedby                         = FieldMask{Size: 2, Offset: 0x000C}
	FieldMaskUnitSummonedby                        = FieldMask{Size: 2, Offset: 0x000E}
	FieldMaskUnitCreatedby                         = FieldMask{Size: 2, Offset: 0x0010}
	FieldMaskUnitTarget                            = FieldMask{Size: 2, Offset: 0x0012}
	FieldMaskUnitChannelObject                     = FieldMask{Size: 2, Offset: 0x0014}
	FieldMaskUnitChannelSpell                      = FieldMask{Size: 1, Offset: 0x0016}
	FieldMaskUnitBytes0                            = FieldMask{Size: 1, Offset: 0x0017}
	FieldMaskUnitHealth                            = FieldMask{Size: 1, Offset: 0x0018}
	FieldMaskUnitPower1                            = FieldMask{Size: 1, Offset: 0x0019}
	FieldMaskUnitPower2                            = FieldMask{Size: 1, Offset: 0x001A}
	FieldMaskUnitPower3                            = FieldMask{Size: 1, Offset: 0x001B}
	FieldMaskUnitPower4                            = FieldMask{Size: 1, Offset: 0x001C}
	FieldMaskUnitPower5                            = FieldMask{Size: 1, Offset: 0x001D}
	FieldMaskUnitPower6                            = FieldMask{Size: 1, Offset: 0x001E}
	FieldMaskUnitPower7                            = FieldMask{Size: 1, Offset: 0x001F}
	FieldMaskUnitMaxhealth                         = FieldMask{Size: 1, Offset: 0x0020}
	FieldMaskUnitMaxpower1                         = FieldMask{Size: 1, Offset: 0x0021}
	FieldMaskUnitMaxpower2                         = FieldMask{Size: 1, Offset: 0x0022}
	FieldMaskUnitMaxpower3                         = FieldMask{Size: 1, Offset: 0x0023}
	FieldMaskUnitMaxpower4                         = FieldMask{Size: 1, Offset: 0x0024}
	FieldMaskUnitMaxpower5                         = FieldMask{Size: 1, Offset: 0x0025}
	FieldMaskUnitMaxpower6                         = FieldMask{Size: 1, Offset: 0x0026}
	FieldMaskUnitMaxpower7                         = FieldMask{Size: 1, Offset: 0x0027}
	FieldMaskUnitPowerRegenFlatModifier            = FieldMask{Size: 7, Offset: 0x0028}
	FieldMaskUnitPowerRegenInterruptedFlatModifier = FieldMask{Size: 7, Offset: 0x002F}
	FieldMaskUnitLevel                             = FieldMask{Size: 1, Offset: 0x0036}
	FieldMaskUnitFactiontemplate                   = FieldMask{Size: 1, Offset: 0x0037}
	FieldMaskUnitVirtualItemSlotId                 = FieldMask{Size: 3, Offset: 0x0038}
	FieldMaskUnitFlags                             = FieldMask{Size: 1, Offset: 0x003B}
	FieldMaskUnitFlags2                            = FieldMask{Size: 1, Offset: 0x003C}
	FieldMaskUnitAurastate                         = FieldMask{Size: 1, Offset: 0x003D}
	FieldMaskUnitBaseattacktime                    = FieldMask{Size: 2, Offset: 0x003E}
	FieldMaskUnitRangedattacktime                  = FieldMask{Size: 1, Offset: 0x0040}
	FieldMaskUnitBoundingradius                    = FieldMask{Size: 1, Offset: 0x0041}
	FieldMaskUnitCombatreach                       = FieldMask{Size: 1, Offset: 0x0042}
	FieldMaskUnitDisplayid                         = FieldMask{Size: 1, Offset: 0x0043}
	FieldMaskUnitNativedisplayid                   = FieldMask{Size: 1, Offset: 0x0044}
	FieldMaskUnitMountdisplayid                    = FieldMask{Size: 1, Offset: 0x0045}
	FieldMaskUnitMindamage                         = FieldMask{Size: 1, Offset: 0x0046}
	FieldMaskUnitMaxdamage                         = FieldMask{Size: 1, Offset: 0x0047}
	FieldMaskUnitMinoffhanddamage                  = FieldMask{Size: 1, Offset: 0x0048}
	FieldMaskUnitMaxoffhanddamage                  = FieldMask{Size: 1, Offset: 0x0049}
	FieldMaskUnitBytes1                            = FieldMask{Size: 1, Offset: 0x004A}
	FieldMaskUnitPetnumber                         = FieldMask{Size: 1, Offset: 0x004B}
	FieldMaskUnitPetNameTimestamp                  = FieldMask{Size: 1, Offset: 0x004C}
	FieldMaskUnitPetexperience                     = FieldMask{Size: 1, Offset: 0x004D}
	FieldMaskUnitPetnextlevelexp                   = FieldMask{Size: 1, Offset: 0x004E}
	FieldMaskUnitDynamicFlags                      = FieldMask{Size: 1, Offset: 0x004F}
	FieldMaskUnitModCastSpeed                      = FieldMask{Size: 1, Offset: 0x0050}
	FieldMaskUnitCreatedBySpell                    = FieldMask{Size: 1, Offset: 0x0051}
	FieldMaskUnitNpcFlags                          = FieldMask{Size: 1, Offset: 0x0052}
	FieldMaskUnitNpcEmotestate                     = FieldMask{Size: 1, Offset: 0x0053}
	FieldMaskUnitStrength                          = FieldMask{Size: 1, Offset: 0x0054}
	FieldMaskUnitAgility                           = FieldMask{Size: 1, Offset: 0x0055}
	FieldMaskUnitStamina                           = FieldMask{Size: 1, Offset: 0x0056}
	FieldMaskUnitIntellect                         = FieldMask{Size: 1, Offset: 0x0057}
	FieldMaskUnitSpirit                            = FieldMask{Size: 1, Offset: 0x0058}
	FieldMaskUnitPosStat0                          = FieldMask{Size: 1, Offset: 0x0059}
	FieldMaskUnitPosStat1                          = FieldMask{Size: 1, Offset: 0x005A}
	FieldMaskUnitPosStat2                          = FieldMask{Size: 1, Offset: 0x005B}
	FieldMaskUnitPosStat3                          = FieldMask{Size: 1, Offset: 0x005C}
	FieldMaskUnitPosStat4                          = FieldMask{Size: 1, Offset: 0x005D}
	FieldMaskUnitNegStat0                          = FieldMask{Size: 1, Offset: 0x005E}
	FieldMaskUnitNegStat1                          = FieldMask{Size: 1, Offset: 0x005F}
	FieldMaskUnitNegStat2                          = FieldMask{Size: 1, Offset: 0x0060}
	FieldMaskUnitNegStat3                          = FieldMask{Size: 1, Offset: 0x0061}
	FieldMaskUnitNegStat4                          = FieldMask{Size: 1, Offset: 0x0062}
	FieldMaskUnitResistances                       = FieldMask{Size: 7, Offset: 0x0063}
	FieldMaskUnitResistanceBuffModsPositive        = FieldMask{Size: 7, Offset: 0x006A}
	FieldMaskUnitResistanceBuffModsNegative        = FieldMask{Size: 7, Offset: 0x0071}
	FieldMaskUnitBaseMana                          = FieldMask{Size: 1, Offset: 0x0078}
	FieldMaskUnitBaseHealth                        = FieldMask{Size: 1, Offset: 0x0079}
	FieldMaskUnitBytes2                            = FieldMask{Size: 1, Offset: 0x007A}
	FieldMaskUnitAttackPower                       = FieldMask{Size: 1, Offset: 0x007B}
	FieldMaskUnitAttackPowerMods                   = FieldMask{Size: 1, Offset: 0x007C}
	FieldMaskUnitAttackPowerMultiplier             = FieldMask{Size: 1, Offset: 0x007D}
	FieldMaskUnitRangedAttackPower                 = FieldMask{Size: 1, Offset: 0x007E}
	FieldMaskUnitRangedAttackPowerMods             = FieldMask{Size: 1, Offset: 0x007F}
	FieldMaskUnitRangedAttackPowerMultiplier       = FieldMask{Size: 1, Offset: 0x0080}
	FieldMaskUnitMinRangedDamage                   = FieldMask{Size: 1, Offset: 0x0081}
	FieldMaskUnitMaxRangedDamage                   = FieldMask{Size: 1, Offset: 0x0082}
	FieldMaskUnitPowerCostModifier                 = FieldMask{Size: 7, Offset: 0x0083}
	FieldMaskUnitPowerCostMultiplier               = FieldMask{Size: 7, Offset: 0x008A}
	FieldMaskUnitMaxHealthModifier                 = FieldMask{Size: 1, Offset: 0x0091}
	FieldMaskUnitHoverHeight                       = FieldMask{Size: 1, Offset: 0x0092}

	FieldMaskPlayerDuelArbiter                 = FieldMask{Size: 2, Offset: 0x0094}
	FieldMaskPlayerFlags                       = FieldMask{Size: 1, Offset: 0x0096}
	FieldMaskPlayerGuildid                     = FieldMask{Size: 1, Offset: 0x0097}
	FieldMaskPlayerGuildrank                   = FieldMask{Size: 1, Offset: 0x0098}
	FieldMaskPlayerFieldBytes                  = FieldMask{Size: 1, Offset: 0x0099}
	FieldMaskPlayerBytes2                      = FieldMask{Size: 1, Offset: 0x009A}
	FieldMaskPlayerBytes3                      = FieldMask{Size: 1, Offset: 0x009B}
	FieldMaskPlayerDuelTeam                    = FieldMask{Size: 1, Offset: 0x009C}
	FieldMaskPlayerGuildTimestamp              = FieldMask{Size: 1, Offset: 0x009D}
	FieldMaskPlayerQuestLog1_1                 = FieldMask{Size: 1, Offset: 0x009E}
	FieldMaskPlayerQuestLog1_2                 = FieldMask{Size: 1, Offset: 0x009F}
	FieldMaskPlayerQuestLog1_3                 = FieldMask{Size: 2, Offset: 0x00A0}
	FieldMaskPlayerQuestLog1_4                 = FieldMask{Size: 1, Offset: 0x00A2}
	FieldMaskPlayerQuestLog2_1                 = FieldMask{Size: 1, Offset: 0x00A3}
	FieldMaskPlayerQuestLog2_2                 = FieldMask{Size: 1, Offset: 0x00A4}
	FieldMaskPlayerQuestLog2_3                 = FieldMask{Size: 2, Offset: 0x00A5}
	FieldMaskPlayerQuestLog2_5                 = FieldMask{Size: 1, Offset: 0x00A7}
	FieldMaskPlayerQuestLog3_1                 = FieldMask{Size: 1, Offset: 0x00A8}
	FieldMaskPlayerQuestLog3_2                 = FieldMask{Size: 1, Offset: 0x00A9}
	FieldMaskPlayerQuestLog3_3                 = FieldMask{Size: 2, Offset: 0x00AA}
	FieldMaskPlayerQuestLog3_5                 = FieldMask{Size: 1, Offset: 0x00AC}
	FieldMaskPlayerQuestLog4_1                 = FieldMask{Size: 1, Offset: 0x00AD}
	FieldMaskPlayerQuestLog4_2                 = FieldMask{Size: 1, Offset: 0x00AE}
	FieldMaskPlayerQuestLog4_3                 = FieldMask{Size: 2, Offset: 0x00AF}
	FieldMaskPlayerQuestLog4_5                 = FieldMask{Size: 1, Offset: 0x00B1}
	FieldMaskPlayerQuestLog5_1                 = FieldMask{Size: 1, Offset: 0x00B2}
	FieldMaskPlayerQuestLog5_2                 = FieldMask{Size: 1, Offset: 0x00B3}
	FieldMaskPlayerQuestLog5_3                 = FieldMask{Size: 2, Offset: 0x00B4}
	FieldMaskPlayerQuestLog5_5                 = FieldMask{Size: 1, Offset: 0x00B6}
	FieldMaskPlayerQuestLog6_1                 = FieldMask{Size: 1, Offset: 0x00B7}
	FieldMaskPlayerQuestLog6_2                 = FieldMask{Size: 1, Offset: 0x00B8}
	FieldMaskPlayerQuestLog6_3                 = FieldMask{Size: 2, Offset: 0x00B9}
	FieldMaskPlayerQuestLog6_5                 = FieldMask{Size: 1, Offset: 0x00BB}
	FieldMaskPlayerQuestLog7_1                 = FieldMask{Size: 1, Offset: 0x00BC}
	FieldMaskPlayerQuestLog7_2                 = FieldMask{Size: 1, Offset: 0x00BD}
	FieldMaskPlayerQuestLog7_3                 = FieldMask{Size: 2, Offset: 0x00BE}
	FieldMaskPlayerQuestLog7_5                 = FieldMask{Size: 1, Offset: 0x00C0}
	FieldMaskPlayerQuestLog8_1                 = FieldMask{Size: 1, Offset: 0x00C1}
	FieldMaskPlayerQuestLog8_2                 = FieldMask{Size: 1, Offset: 0x00C2}
	FieldMaskPlayerQuestLog8_3                 = FieldMask{Size: 2, Offset: 0x00C3}
	FieldMaskPlayerQuestLog8_5                 = FieldMask{Size: 1, Offset: 0x00C5}
	FieldMaskPlayerQuestLog9_1                 = FieldMask{Size: 1, Offset: 0x00C6}
	FieldMaskPlayerQuestLog9_2                 = FieldMask{Size: 1, Offset: 0x00C7}
	FieldMaskPlayerQuestLog9_3                 = FieldMask{Size: 2, Offset: 0x00C8}
	FieldMaskPlayerQuestLog9_5                 = FieldMask{Size: 1, Offset: 0x00CA}
	FieldMaskPlayerQuestLog10_1                = FieldMask{Size: 1, Offset: 0x00CB}
	FieldMaskPlayerQuestLog10_2                = FieldMask{Size: 1, Offset: 0x00CC}
	FieldMaskPlayerQuestLog10_3                = FieldMask{Size: 2, Offset: 0x00CD}
	FieldMaskPlayerQuestLog10_5                = FieldMask{Size: 1, Offset: 0x00CF}
	FieldMaskPlayerQuestLog11_1                = FieldMask{Size: 1, Offset: 0x00D0}
	FieldMaskPlayerQuestLog11_2                = FieldMask{Size: 1, Offset: 0x00D1}
	FieldMaskPlayerQuestLog11_3                = FieldMask{Size: 2, Offset: 0x00D2}
	FieldMaskPlayerQuestLog11_5                = FieldMask{Size: 1, Offset: 0x00D4}
	FieldMaskPlayerQuestLog12_1                = FieldMask{Size: 1, Offset: 0x00D5}
	FieldMaskPlayerQuestLog12_2                = FieldMask{Size: 1, Offset: 0x00D6}
	FieldMaskPlayerQuestLog12_3                = FieldMask{Size: 2, Offset: 0x00D7}
	FieldMaskPlayerQuestLog12_5                = FieldMask{Size: 1, Offset: 0x00D9}
	FieldMaskPlayerQuestLog13_1                = FieldMask{Size: 1, Offset: 0x00DA}
	FieldMaskPlayerQuestLog13_2                = FieldMask{Size: 1, Offset: 0x00DB}
	FieldMaskPlayerQuestLog13_3                = FieldMask{Size: 2, Offset: 0x00DC}
	FieldMaskPlayerQuestLog13_5                = FieldMask{Size: 1, Offset: 0x00DE}
	FieldMaskPlayerQuestLog14_1                = FieldMask{Size: 1, Offset: 0x00DF}
	FieldMaskPlayerQuestLog14_2                = FieldMask{Size: 1, Offset: 0x00E0}
	FieldMaskPlayerQuestLog14_3                = FieldMask{Size: 2, Offset: 0x00E1}
	FieldMaskPlayerQuestLog14_5                = FieldMask{Size: 1, Offset: 0x00E3}
	FieldMaskPlayerQuestLog15_1                = FieldMask{Size: 1, Offset: 0x00E4}
	FieldMaskPlayerQuestLog15_2                = FieldMask{Size: 1, Offset: 0x00E5}
	FieldMaskPlayerQuestLog15_3                = FieldMask{Size: 2, Offset: 0x00E6}
	FieldMaskPlayerQuestLog15_5                = FieldMask{Size: 1, Offset: 0x00E8}
	FieldMaskPlayerQuestLog16_1                = FieldMask{Size: 1, Offset: 0x00E9}
	FieldMaskPlayerQuestLog16_2                = FieldMask{Size: 1, Offset: 0x00EA}
	FieldMaskPlayerQuestLog16_3                = FieldMask{Size: 2, Offset: 0x00EB}
	FieldMaskPlayerQuestLog16_5                = FieldMask{Size: 1, Offset: 0x00ED}
	FieldMaskPlayerQuestLog17_1                = FieldMask{Size: 1, Offset: 0x00EE}
	FieldMaskPlayerQuestLog17_2                = FieldMask{Size: 1, Offset: 0x00EF}
	FieldMaskPlayerQuestLog17_3                = FieldMask{Size: 2, Offset: 0x00F0}
	FieldMaskPlayerQuestLog17_5                = FieldMask{Size: 1, Offset: 0x00F2}
	FieldMaskPlayerQuestLog18_1                = FieldMask{Size: 1, Offset: 0x00F3}
	FieldMaskPlayerQuestLog18_2                = FieldMask{Size: 1, Offset: 0x00F4}
	FieldMaskPlayerQuestLog18_3                = FieldMask{Size: 2, Offset: 0x00F5}
	FieldMaskPlayerQuestLog18_5                = FieldMask{Size: 1, Offset: 0x00F7}
	FieldMaskPlayerQuestLog19_1                = FieldMask{Size: 1, Offset: 0x00F8}
	FieldMaskPlayerQuestLog19_2                = FieldMask{Size: 1, Offset: 0x00F9}
	FieldMaskPlayerQuestLog19_3                = FieldMask{Size: 2, Offset: 0x00FA}
	FieldMaskPlayerQuestLog19_5                = FieldMask{Size: 1, Offset: 0x00FC}
	FieldMaskPlayerQuestLog20_1                = FieldMask{Size: 1, Offset: 0x00FD}
	FieldMaskPlayerQuestLog20_2                = FieldMask{Size: 1, Offset: 0x00FE}
	FieldMaskPlayerQuestLog20_3                = FieldMask{Size: 2, Offset: 0x00FF}
	FieldMaskPlayerQuestLog20_5                = FieldMask{Size: 1, Offset: 0x0101}
	FieldMaskPlayerQuestLog21_1                = FieldMask{Size: 1, Offset: 0x0102}
	FieldMaskPlayerQuestLog21_2                = FieldMask{Size: 1, Offset: 0x0103}
	FieldMaskPlayerQuestLog21_3                = FieldMask{Size: 2, Offset: 0x0104}
	FieldMaskPlayerQuestLog21_5                = FieldMask{Size: 1, Offset: 0x0106}
	FieldMaskPlayerQuestLog22_1                = FieldMask{Size: 1, Offset: 0x0107}
	FieldMaskPlayerQuestLog22_2                = FieldMask{Size: 1, Offset: 0x0108}
	FieldMaskPlayerQuestLog22_3                = FieldMask{Size: 2, Offset: 0x0109}
	FieldMaskPlayerQuestLog22_5                = FieldMask{Size: 1, Offset: 0x010B}
	FieldMaskPlayerQuestLog23_1                = FieldMask{Size: 1, Offset: 0x010C}
	FieldMaskPlayerQuestLog23_2                = FieldMask{Size: 1, Offset: 0x010D}
	FieldMaskPlayerQuestLog23_3                = FieldMask{Size: 2, Offset: 0x010E}
	FieldMaskPlayerQuestLog23_5                = FieldMask{Size: 1, Offset: 0x0110}
	FieldMaskPlayerQuestLog24_1                = FieldMask{Size: 1, Offset: 0x0111}
	FieldMaskPlayerQuestLog24_2                = FieldMask{Size: 1, Offset: 0x0112}
	FieldMaskPlayerQuestLog24_3                = FieldMask{Size: 2, Offset: 0x0113}
	FieldMaskPlayerQuestLog24_5                = FieldMask{Size: 1, Offset: 0x0115}
	FieldMaskPlayerQuestLog25_1                = FieldMask{Size: 1, Offset: 0x0116}
	FieldMaskPlayerQuestLog25_2                = FieldMask{Size: 1, Offset: 0x0117}
	FieldMaskPlayerQuestLog25_3                = FieldMask{Size: 2, Offset: 0x0118}
	FieldMaskPlayerQuestLog25_5                = FieldMask{Size: 1, Offset: 0x011A}
	FieldMaskPlayerVisibleItem                 = FieldMask{Size: 3, Offset: 0x011B}
	FieldMaskPlayerChosenTitle                 = FieldMask{Size: 1, Offset: 0x0141}
	FieldMaskPlayerFakeInebriation             = FieldMask{Size: 1, Offset: 0x0142}
	FieldMaskPlayerFieldInv                    = FieldMask{Size: 3, Offset: 0x0144}
	FieldMaskPlayerFarsight                    = FieldMask{Size: 2, Offset: 0x0270}
	FieldMaskPlayerKnownTitles                 = FieldMask{Size: 2, Offset: 0x0272}
	FieldMaskPlayerKnownTitles1                = FieldMask{Size: 2, Offset: 0x0274}
	FieldMaskPlayerKnownTitles2                = FieldMask{Size: 2, Offset: 0x0276}
	FieldMaskPlayerKnownCurrencies             = FieldMask{Size: 2, Offset: 0x0278}
	FieldMaskPlayerXp                          = FieldMask{Size: 1, Offset: 0x027A}
	FieldMaskPlayerNextLevelXp                 = FieldMask{Size: 1, Offset: 0x027B}
	FieldMaskPlayerSkillInfo                   = FieldMask{Size: 3, Offset: 0x027C}
	FieldMaskPlayerCharacterPoints1            = FieldMask{Size: 1, Offset: 0x03FC}
	FieldMaskPlayerCharacterPoints2            = FieldMask{Size: 1, Offset: 0x03FD}
	FieldMaskPlayerTrackCreatures              = FieldMask{Size: 1, Offset: 0x03FE}
	FieldMaskPlayerTrackResources              = FieldMask{Size: 1, Offset: 0x03FF}
	FieldMaskPlayerBlockPercentage             = FieldMask{Size: 1, Offset: 0x0400}
	FieldMaskPlayerDodgePercentage             = FieldMask{Size: 1, Offset: 0x0401}
	FieldMaskPlayerParryPercentage             = FieldMask{Size: 1, Offset: 0x0402}
	FieldMaskPlayerExpertise                   = FieldMask{Size: 1, Offset: 0x0403}
	FieldMaskPlayerOffhandExpertise            = FieldMask{Size: 1, Offset: 0x0404}
	FieldMaskPlayerCritPercentage              = FieldMask{Size: 1, Offset: 0x0405}
	FieldMaskPlayerRangedCritPercentage        = FieldMask{Size: 1, Offset: 0x0406}
	FieldMaskPlayerOffhandCritPercentage       = FieldMask{Size: 1, Offset: 0x0407}
	FieldMaskPlayerSpellCritPercentage1        = FieldMask{Size: 7, Offset: 0x0408}
	FieldMaskPlayerShieldBlock                 = FieldMask{Size: 1, Offset: 0x040F}
	FieldMaskPlayerShieldBlockCritPercentage   = FieldMask{Size: 1, Offset: 0x0410}
	FieldMaskPlayerExploredZones1              = FieldMask{Size: 1, Offset: 0x0411}
	FieldMaskPlayerRestStateExperience         = FieldMask{Size: 1, Offset: 0x0491}
	FieldMaskPlayerCoinage                     = FieldMask{Size: 1, Offset: 0x0492}
	FieldMaskPlayerModDamageDonePos            = FieldMask{Size: 7, Offset: 0x0493}
	FieldMaskPlayerModDamageDoneNeg            = FieldMask{Size: 7, Offset: 0x049A}
	FieldMaskPlayerModDamageDonePct            = FieldMask{Size: 7, Offset: 0x04A1}
	FieldMaskPlayerModHealingDonePos           = FieldMask{Size: 1, Offset: 0x04A8}
	FieldMaskPlayerModHealingPct               = FieldMask{Size: 1, Offset: 0x04A9}
	FieldMaskPlayerModHealingDonePct           = FieldMask{Size: 1, Offset: 0x04AA}
	FieldMaskPlayerModTargetResistance         = FieldMask{Size: 1, Offset: 0x04AB}
	FieldMaskPlayerModTargetPhysicalResistance = FieldMask{Size: 1, Offset: 0x04AC}
	FieldMaskPlayerFeatures                    = FieldMask{Size: 1, Offset: 0x04AD}
	FieldMaskPlayerAmmoId                      = FieldMask{Size: 1, Offset: 0x04AE}
	FieldMaskPlayerSelfResSpell                = FieldMask{Size: 1, Offset: 0x04AF}
	FieldMaskPlayerPvpMedals                   = FieldMask{Size: 1, Offset: 0x04B0}
	FieldMaskPlayerBuybackPrice1               = FieldMask{Size: 1, Offset: 0x04B1}
	FieldMaskPlayerBuybackTimestamp1           = FieldMask{Size: 1, Offset: 0x04BD}
	FieldMaskPlayerKills                       = FieldMask{Size: 1, Offset: 0x04C9}
	FieldMaskPlayerTodayContribution           = FieldMask{Size: 1, Offset: 0x04CA}
	FieldMaskPlayerYesterdayContribution       = FieldMask{Size: 1, Offset: 0x04CB}
	FieldMaskPlayerLifetimeHonorableKills      = FieldMask{Size: 1, Offset: 0x04CC}
	FieldMaskPlayerBytes4                      = FieldMask{Size: 1, Offset: 0x04CD}
	FieldMaskPlayerWatchedFactionIndex         = FieldMask{Size: 1, Offset: 0x04CE}
	FieldMaskPlayerCombatRating1               = FieldMask{Size: 2, Offset: 0x04CF}
	FieldMaskPlayerArenaTeamInfo11             = FieldMask{Size: 2, Offset: 0x04E8}
	FieldMaskPlayerHonorCurrency               = FieldMask{Size: 1, Offset: 0x04FD}
	FieldMaskPlayerArenaCurrency               = FieldMask{Size: 1, Offset: 0x04FE}
	FieldMaskPlayerMaxLevel                    = FieldMask{Size: 1, Offset: 0x04FF}
	FieldMaskPlayerDailyQuests1                = FieldMask{Size: 2, Offset: 0x0500}
	FieldMaskPlayerRuneRegen1                  = FieldMask{Size: 4, Offset: 0x0519}
	FieldMaskPlayerNoReagentCost1              = FieldMask{Size: 3, Offset: 0x051D}
	FieldMaskPlayerGlyphSlots1                 = FieldMask{Size: 6, Offset: 0x0520}
	FieldMaskPlayerGlyphs1                     = FieldMask{Size: 6, Offset: 0x0526}
	FieldMaskPlayerGlyphsEnabled               = FieldMask{Size: 1, Offset: 0x052C}
	FieldMaskPlayerPetSpellPower               = FieldMask{Size: 1, Offset: 0x052D}

	FieldMaskGameObjectDisplayid      = FieldMask{Size: 1, Offset: 0x0008}
	FieldMaskGameObjectFlags          = FieldMask{Size: 1, Offset: 0x0009}
	FieldMaskGameObjectParentrotation = FieldMask{Size: 4, Offset: 0x000A}
	FieldMaskGameObjectDynamic        = FieldMask{Size: 1, Offset: 0x000E}
	FieldMaskGameObjectFaction        = FieldMask{Size: 1, Offset: 0x000F}
	FieldMaskGameObjectLevel          = FieldMask{Size: 1, Offset: 0x0010}
	FieldMaskGameObjectBytes1         = FieldMask{Size: 1, Offset: 0x0011}

	FieldMaskDynamicObjectCaster   = FieldMask{Size: 2, Offset: 0x0006}
	FieldMaskDynamicObjectBytes    = FieldMask{Size: 1, Offset: 0x0008}
	FieldMaskDynamicObjectSpellid  = FieldMask{Size: 1, Offset: 0x0009}
	FieldMaskDynamicObjectRadius   = FieldMask{Size: 1, Offset: 0x000A}
	FieldMaskDynamicObjectCasttime = FieldMask{Size: 1, Offset: 0x000B}

	FieldMaskCorpseOwner        = FieldMask{Size: 2, Offset: 0x0006}
	FieldMaskCorpseParty        = FieldMask{Size: 2, Offset: 0x0008}
	FieldMaskCorpseDisplayId    = FieldMask{Size: 1, Offset: 0x000A}
	FieldMaskCorpseItem         = FieldMask{Size: 1, Offset: 0x000B}
	FieldMaskCorpseBytes1       = FieldMask{Size: 1, Offset: 0x001E}
	FieldMaskCorpseBytes2       = FieldMask{Size: 1, Offset: 0x001F}
	FieldMaskCorpseGuild        = FieldMask{Size: 1, Offset: 0x0020}
	FieldMaskCorpseFlags        = FieldMask{Size: 1, Offset: 0x0021}
	FieldMaskCorpseDynamicFlags = FieldMask{Size: 1, Offset: 0x0022}
)
