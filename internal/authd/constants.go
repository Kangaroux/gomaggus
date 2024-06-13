package authd

type Opcode byte

const (
	OpLoginChallenge     Opcode = 0x0
	OpLoginProof         Opcode = 0x1
	OpReconnectChallenge Opcode = 0x2
	OpReconnectProof     Opcode = 0x3
	OpRealmList          Opcode = 0x10
)

type ErrorCode byte

const (
	CodeSuccess            ErrorCode = 0x0
	CodeFailUnknownAccount ErrorCode = 0x4
)

type RealmFlag uint8

const (
	RealmFlagNone         RealmFlag = 0x0
	RealmFlagInvalid      RealmFlag = 0x1 // Realm is greyed out and can't be selected
	RealmFlagOffline      RealmFlag = 0x2 // Population: "Offline" and can't be selected
	RealmFlagSpecifyBuild RealmFlag = 0x4 // Includes version in realm name
	RealmFlagUnknown1     RealmFlag = 0x8
	RealmFlagUnknown2     RealmFlag = 0x10
	RealmFlagNewPlayers   RealmFlag = 0x20 // Population: "New Players" in blue text
	RealmFlagNewServer    RealmFlag = 0x40 // Population: "New" in green text
	RealmFlagFull         RealmFlag = 0x80 // Population: "Full" in red text
)

const (
	ReconnectDataLen = 16
)
