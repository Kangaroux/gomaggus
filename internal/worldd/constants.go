package worldd

// Opcodes sent by the server
const (
	OP_AUTH_CHALLENGE uint16 = 0x1EC
	OP_AUTH_RESPONSE  uint16 = 0x1EE
)

// Opcodes sent by the client
const (
	OP_AUTH_SESSION uint32 = 0x1ED
)

const (
	RespCodeResponseSuccess                                byte = 0x00
	RespCodeResponseFailure                                byte = 0x01
	RespCodeResponseCancelled                              byte = 0x02
	RespCodeResponseDisconnected                           byte = 0x03
	RespCodeResponseFailedToConnect                        byte = 0x04
	RespCodeResponseConnected                              byte = 0x05
	RespCodeResponseVersionMismatch                        byte = 0x06
	RespCodeCStatusConnecting                              byte = 0x07
	RespCodeCStatusNegotiatingSecurity                     byte = 0x08
	RespCodeCStatusNegotiationComplete                     byte = 0x09
	RespCodeCStatusNegotiationFailed                       byte = 0x0A
	RespCodeCStatusAuthenticating                          byte = 0x0B
	RespCodeAuthOk                                         byte = 0x0C
	RespCodeAuthFailed                                     byte = 0x0D
	RespCodeAuthReject                                     byte = 0x0E
	RespCodeAuthBadServerProof                             byte = 0x0F
	RespCodeAuthUnavailable                                byte = 0x10
	RespCodeAuthSystemError                                byte = 0x11
	RespCodeAuthBillingError                               byte = 0x12
	RespCodeAuthBillingExpired                             byte = 0x13
	RespCodeAuthVersionMismatch                            byte = 0x14
	RespCodeAuthUnknownAccount                             byte = 0x15
	RespCodeAuthIncorrectPassword                          byte = 0x16
	RespCodeAuthSessionExpired                             byte = 0x17
	RespCodeAuthServerShuttingDown                         byte = 0x18
	RespCodeAuthAlreadyLoggingIn                           byte = 0x19
	RespCodeAuthLoginServerNotFound                        byte = 0x1A
	RespCodeAuthWaitQueue                                  byte = 0x1B
	RespCodeAuthBanned                                     byte = 0x1C
	RespCodeAuthAlreadyOnline                              byte = 0x1D
	RespCodeAuthNoTime                                     byte = 0x1E
	RespCodeAuthDbBusy                                     byte = 0x1F
	RespCodeAuthSuspended                                  byte = 0x20
	RespCodeAuthParentalControl                            byte = 0x21
	RespCodeAuthLockedEnforced                             byte = 0x22
	RespCodeRealmListInProgress                            byte = 0x23
	RespCodeRealmListSuccess                               byte = 0x24
	RespCodeRealmListFailed                                byte = 0x25
	RespCodeRealmListInvalid                               byte = 0x26
	RespCodeRealmListRealmNotFound                         byte = 0x27
	RespCodeAccountCreateInProgress                        byte = 0x28
	RespCodeAccountCreateSuccess                           byte = 0x29
	RespCodeAccountCreateFailed                            byte = 0x2A
	RespCodeCharListRetrieving                             byte = 0x2B
	RespCodeCharListRetrieved                              byte = 0x2C
	RespCodeCharListFailed                                 byte = 0x2D
	RespCodeCharCreateInProgress                           byte = 0x2E
	RespCodeCharCreateSuccess                              byte = 0x2F
	RespCodeCharCreateError                                byte = 0x30
	RespCodeCharCreateFailed                               byte = 0x31
	RespCodeCharCreateNameInUse                            byte = 0x32
	RespCodeCharCreateDisabled                             byte = 0x33
	RespCodeCharCreatePvpTeamsViolation                    byte = 0x34
	RespCodeCharCreateServerLimit                          byte = 0x35
	RespCodeCharCreateAccountLimit                         byte = 0x36
	RespCodeCharCreateServerQueue                          byte = 0x37
	RespCodeCharCreateOnlyExisting                         byte = 0x38
	RespCodeCharCreateExpansion                            byte = 0x39
	RespCodeCharCreateExpansionClass                       byte = 0x3A
	RespCodeCharCreateLevelRequirement                     byte = 0x3B
	RespCodeCharCreateUniqueClassLimit                     byte = 0x3C
	RespCodeCharCreateCharacterInGuild                     byte = 0x3D
	RespCodeCharCreateRestrictedRaceclass                  byte = 0x3E
	RespCodeCharCreateCharacterChooseRace                  byte = 0x3F
	RespCodeCharCreateCharacterArenaLeader                 byte = 0x40
	RespCodeCharCreateCharacterDeleteMail                  byte = 0x41
	RespCodeCharCreateCharacterSwapFaction                 byte = 0x42
	RespCodeCharCreateCharacterRaceOnly                    byte = 0x43
	RespCodeCharCreateCharacterGoldLimit                   byte = 0x44
	RespCodeCharCreateForceLogin                           byte = 0x45
	RespCodeCharDeleteInProgress                           byte = 0x46
	RespCodeCharDeleteSuccess                              byte = 0x47
	RespCodeCharDeleteFailed                               byte = 0x48
	RespCodeCharDeleteFailedLockedForTransfer              byte = 0x49
	RespCodeCharDeleteFailedGuildLeader                    byte = 0x4A
	RespCodeCharDeleteFailedArenaCaptain                   byte = 0x4B
	RespCodeCharLoginInProgress                            byte = 0x4C
	RespCodeCharLoginSuccess                               byte = 0x4D
	RespCodeCharLoginNoWorld                               byte = 0x4E
	RespCodeCharLoginDuplicateCharacter                    byte = 0x4F
	RespCodeCharLoginNoInstances                           byte = 0x50
	RespCodeCharLoginFailed                                byte = 0x51
	RespCodeCharLoginDisabled                              byte = 0x52
	RespCodeCharLoginNoCharacter                           byte = 0x53
	RespCodeCharLoginLockedForTransfer                     byte = 0x54
	RespCodeCharLoginLockedByBilling                       byte = 0x55
	RespCodeCharLoginLockedByMobileAh                      byte = 0x56
	RespCodeCharNameSuccess                                byte = 0x57
	RespCodeCharNameFailure                                byte = 0x58
	RespCodeCharNameNoName                                 byte = 0x59
	RespCodeCharNameTooShort                               byte = 0x5A
	RespCodeCharNameTooLong                                byte = 0x5B
	RespCodeCharNameInvalidCharacter                       byte = 0x5C
	RespCodeCharNameMixedLanguages                         byte = 0x5D
	RespCodeCharNameProfane                                byte = 0x5E
	RespCodeCharNameReserved                               byte = 0x5F
	RespCodeCharNameInvalidApostrophe                      byte = 0x60
	RespCodeCharNameMultipleApostrophes                    byte = 0x61
	RespCodeCharNameThreeConsecutive                       byte = 0x62
	RespCodeCharNameInvalidSpace                           byte = 0x63
	RespCodeCharNameConsecutiveSpaces                      byte = 0x64
	RespCodeCharNameRussianConsecutiveSilentCharacters     byte = 0x65
	RespCodeCharNameRussianSilentCharacterAtBeginningOrEnd byte = 0x66
	RespCodeCharNameDeclensionDoesntMatchBaseName          byte = 0x67
)

const (
	ExpansionVanilla byte = 0x0
	ExpansionTbc     byte = 0x1
	ExpansionWrath   byte = 0x2
)
