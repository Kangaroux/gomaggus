package realmd

type ServerOpcode uint16

const (
	OpServerAuthChallenge        ServerOpcode = 0x1EC
	OpServerAuthResponse         ServerOpcode = 0x1EE
	OpServerPong                 ServerOpcode = 0x1DD
	OpServerAccountDataTimes     ServerOpcode = 0x209
	OpServerCharEnum             ServerOpcode = 0x3B
	OpServerRealmSplit           ServerOpcode = 0x38B
	OpServerCharCreate           ServerOpcode = 0x3A
	OpServerCharDelete           ServerOpcode = 0x3C
	OpServerCharLoginFailed      ServerOpcode = 0x41
	OpServerCharLoginVerifyWorld ServerOpcode = 0x236
	OpServerUpdateObject         ServerOpcode = 0xA9
	OpServerTutorialFlags        ServerOpcode = 0xFD
	OpServerSystemFeatures       ServerOpcode = 0x3C9
	OpServerHearthLocation       ServerOpcode = 0x155 // SMSG_BINDPOINTUPDATE
	OpServerPlayCinematic        ServerOpcode = 0xFA  // SMSG_TRIGGER_CINEMATIC
)

type ClientOpcode uint32

const (
	OpClientAuthSession              ClientOpcode = 0x1ED
	OpClientRealmSplit               ClientOpcode = 0x38C
	OpClientPing                     ClientOpcode = 0x1DC
	OpClientCharList                 ClientOpcode = 0x37
	OpClientCharCreate               ClientOpcode = 0x36
	OpClientCharDelete               ClientOpcode = 0x38
	OpClientPlayerLogin              ClientOpcode = 0x3D
	OpClientReadyForAccountDataTimes ClientOpcode = 0x4FF
)
