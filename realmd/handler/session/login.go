package session

import (
	"bytes"
	"database/sql"
	"math"
	"time"

	"github.com/kangaroux/gomaggus/realmd"
	"github.com/kangaroux/gomaggus/realmd/objupdate"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmsg_player_login.html
type loginRequest struct {
	CharacterId uint64
}

// https: gtker.com/wow_messages/docs/smsg_character_login_failed.html#client-version-335
type loginFailed struct {
	ResponseCode realmd.ResponseCode
}

func LoginHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	req := loginRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	char, err := svc.Characters.Get(uint32(req.CharacterId))
	if err != nil {
		return err
	}

	// Notify the client something isn't right and close the connection
	if char == nil || char.AccountId != client.Account.Id || char.RealmId != client.Realm.Id {
		resp := loginFailed{ResponseCode: realmd.RespCodeCharLoginFailed}

		if err := client.SendPacket(realmd.OpServerCharLoginFailed, &resp); err != nil {
			return err
		}

		client.Log.Warn().Uint64("char", req.CharacterId).Msg("client tried logging in as invalid character")

		// TODO: set client state as invalid
		return &realmd.ErrKickClient{Reason: "invalid login"}
	}

	client.Character = char

	if err := sendVerifyWorld(client); err != nil {
		return err
	}
	if err := sendTutorialFlags(client); err != nil {
		return err
	}
	if err := sendSystemFeatures(client); err != nil {
		return err
	}
	if err := sendHearthLocation(client); err != nil {
		return err
	}
	if err := sendIntroCinematic(client); err != nil {
		return err
	}
	if err := sendSpawnPlayer(client); err != nil {
		return err
	}

	client.Character.LastLogin = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if _, err := svc.Characters.Update(client.Character); err != nil {
		return err
	}

	client.Log.Info().Str("char", client.Character.String()).Msg("player login")

	return nil
}

// https://gtker.com/wow_messages/docs/smsg_login_verify_world.html
type verifyWorldResponse struct {
	Map      uint32
	Position realmd.Vector4
}

func sendVerifyWorld(client *realmd.Client) error {
	resp := verifyWorldResponse{
		Map: 0x0, // Eastern kingdoms
		Position: realmd.Vector4{
			X:        float32(-8949.95),
			Y:        float32(-132.493),
			Z:        float32(83.5312),
			Rotation: float32(0),
		},
	}

	return client.SendPacket(realmd.OpServerCharLoginVerifyWorld, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_tutorial_flags.html
type tutorialFlags struct {
	Flags []byte `binary:"[32]byte"`
}

func sendTutorialFlags(client *realmd.Client) error {
	// Enable all tutorial flags
	resp := tutorialFlags{Flags: bytes.Repeat([]byte{0xFF}, 32)}
	return client.SendPacket(realmd.OpServerTutorialFlags, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_feature_system_status.html#client-version-335
type systemFeatures struct {
	// player reporting?
	ComplaintStatus uint8 // 0=disabled, 1=enabled (no auto ignore), 2=enabled (auto ignore)
	VoipEnabled     bool
}

func sendSystemFeatures(client *realmd.Client) error {
	resp := systemFeatures{
		ComplaintStatus: 0x2,
		VoipEnabled:     false,
	}
	return client.SendPacket(realmd.OpServerSystemFeatures, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_bindpointupdate.html#client-version-335
type hearthLocation struct {
	Position realmd.Vector3
	Map      uint32
	Area     uint32
}

func sendHearthLocation(client *realmd.Client) error {
	resp := hearthLocation{
		Position: realmd.Vector3{
			X: float32(-8949.95),
			Y: float32(-132.493),
			Z: float32(83.5312),
		},
		Map:  0x0, // Eastern kingdoms
		Area: 0xB, // Elwynn forest
	}
	return client.SendPacket(realmd.OpServerHearthLocation, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_trigger_cinematic.html
type playCinematic struct {
	CinematicId uint32
}

func sendIntroCinematic(client *realmd.Client) error {
	// Only play the cinematic on first login
	if client.Character.LastLogin.Valid {
		return nil
	}

	resp := playCinematic{CinematicId: 81} // 81 = human
	return client.SendPacket(realmd.OpServerPlayCinematic, &resp)
}

// https://gtker.com/wow_messages/docs/smsg_update_object.html#client-version-335
// TODO: need a builder for these packets
func sendSpawnPlayer(client *realmd.Client) error {
	char := client.Character
	inner := bytes.Buffer{}
	inner.Write([]byte{1, 0, 0, 0}) // number of objects

	// nested object start
	inner.WriteByte(byte(objupdate.UpdateTypeCreateObject))
	inner.Write(realmd.PackGuid(uint64(char.Id)))
	inner.WriteByte(byte(objupdate.ObjectTypePlayer))

	movement := objupdate.MovementValues{}
	movement.Self()

	living := movement.Living()
	living.Data(&objupdate.LivingData{
		Timestamp: uint32(time.Now().Unix()),
		PositionRotation: realmd.Vector4{
			X:        float32(-8949.95),
			Y:        float32(-132.493),
			Z:        float32(83.5312),
			Rotation: float32(0),
		},
		FallTime:           0,
		WalkSpeed:          1,
		RunSpeed:           70,
		ReverseSpeed:       4.5,
		SwimSpeed:          0,
		SwimReverseSpeed:   0,
		FlightSpeed:        0,
		FlightReverseSpeed: 0,
		TurnRate:           math.Pi,
		PitchRate:          0,
	})

	inner.Write(movement.Bytes())

	values := objupdate.Values{}
	obj := values.Object()
	obj.Guid(realmd.Guid(char.Id))
	obj.Type(objupdate.ObjectTypeObject, objupdate.ObjectTypeUnit, objupdate.ObjectTypePlayer)
	obj.ScaleX(1)

	unit := values.Unit()
	unit.RaceClassGenderPower(char.Race, char.Class, char.Gender, realmd.PowerTypeForClass(char.Class))
	unit.Health(100)
	unit.MaxHealth(100)
	unit.Level(1)
	unit.Faction(char.Race)
	unit.DisplayModel(0x4D0C)       // human female
	unit.NativeDisplayModel(0x4D0C) // human female

	inner.Write(values.Bytes())

	return client.SendPacketBytes(realmd.OpServerUpdateObject, inner.Bytes())
}
