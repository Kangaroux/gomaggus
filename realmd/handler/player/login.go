package player

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/realmd"
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

func LoginHandler(svc *realmd.Service, client *realmd.Client, data *realmd.ClientPacket) error {
	log.Println("start character login")

	p := loginRequest{}

	if _, err := binarystruct.Unmarshal(data.Payload, binarystruct.LittleEndian, &p); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(p.CharacterId))
	if err != nil {
		return err
	}

	// Notify the client something isn't right and close the connection
	if char == nil || char.AccountId != client.Account.Id || char.RealmId != client.Realm.Id {
		resp := loginFailed{ResponseCode: realmd.RespCodeCharLoginFailed}

		if err := client.SendPacket(realmd.OpServerCharLoginFailed, &resp); err != nil {
			return err
		}

		// TODO: set client state as invalid / close connection
		return nil
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
	if _, err := svc.Chars.Update(client.Character); err != nil {
		return err
	}

	log.Println("finished character login")

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

	if err := client.SendPacket(realmd.OpServerCharLoginVerifyWorld, &resp); err != nil {
		return err
	}

	log.Println("sent verify world")
	return nil
}

// https://gtker.com/wow_messages/docs/smsg_tutorial_flags.html
type tutorialFlags struct {
	Flags []byte `binary:"[32]byte"`
}

func sendTutorialFlags(client *realmd.Client) error {
	// Enable all tutorial flags
	resp := tutorialFlags{Flags: bytes.Repeat([]byte{0xFF}, 32)}
	if err := client.SendPacket(realmd.OpServerTutorialFlags, &resp); err != nil {
		return err
	}

	log.Println("sent tutorial flags")
	return nil
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
	if err := client.SendPacket(realmd.OpServerSystemFeatures, &resp); err != nil {
		return err
	}

	log.Println("sent system features")
	return nil
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
	if err := client.SendPacket(realmd.OpServerHearthLocation, &resp); err != nil {
		return err
	}

	log.Println("sent hearth location")
	return nil
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
	if err := client.SendPacket(realmd.OpServerPlayCinematic, &resp); err != nil {
		return err
	}

	log.Println("sent intro cinematic")
	return nil
}

// https://gtker.com/wow_messages/docs/smsg_update_object.html#client-version-335
// TODO: need a builder for these packets
func sendSpawnPlayer(client *realmd.Client) error {
	char := client.Character
	inner := bytes.Buffer{}
	inner.Write([]byte{1, 0, 0, 0}) // number of objects

	// nested object start
	inner.WriteByte(byte(realmd.UpdateTypeCreateObject2)) // update type: CREATE_OBJECT2
	inner.Write(realmd.PackGuid(uint64(char.Id)))         // packed guid
	inner.WriteByte(byte(realmd.ObjectTypePlayer))

	// movement block start
	// inner.WriteByte()
	binary.Write(&inner, binary.LittleEndian, realmd.UpdateFlagSelf|realmd.UpdateFlagLiving)
	inner.Write([]byte{0, 0, 0, 0, 0, 0})                        // movement flags
	inner.Write([]byte{0, 0, 0, 0})                              // timestamp
	binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // x
	binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // y
	binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // z
	binary.Write(&inner, binary.LittleEndian, float32(0))        // orientation
	inner.Write([]byte{0, 0, 0, 0})                              // fall time

	binary.Write(&inner, binary.LittleEndian, float32(1))       // walk speed
	binary.Write(&inner, binary.LittleEndian, float32(70))      // run speed
	binary.Write(&inner, binary.LittleEndian, float32(4.5))     // reverse speed
	binary.Write(&inner, binary.LittleEndian, float32(0))       // swim speed
	binary.Write(&inner, binary.LittleEndian, float32(0))       // swim reverse speed
	binary.Write(&inner, binary.LittleEndian, float32(0))       // flight speed
	binary.Write(&inner, binary.LittleEndian, float32(0))       // flight reverse speed
	binary.Write(&inner, binary.LittleEndian, float32(3.14159)) // turn speed
	binary.Write(&inner, binary.LittleEndian, float32(0))       // pitch rate
	// movement block end

	// field mask start
	updateMask := realmd.NewUpdateMask()
	valuesBuf := bytes.Buffer{}

	// Without this, client gets stuck on loading screen and floods server with 0x2CE opcode
	updateMask.SetFieldMask(realmd.FieldMaskObjectGuid)
	binary.Write(&valuesBuf, binary.LittleEndian, uint32(char.Id)) // low guid
	binary.Write(&valuesBuf, binary.LittleEndian, uint32(0))       // high guid

	// Character seems to load fine without this
	updateMask.SetFieldMask(realmd.FieldMaskObjectType)
	binary.Write(&valuesBuf, binary.LittleEndian, uint32(1<<realmd.ObjectTypeObject|1<<realmd.ObjectTypeUnit|1<<realmd.ObjectTypePlayer))

	// Without this, character model scale is zero and camera starts in first person
	updateMask.SetFieldMask(realmd.FieldMaskObjectScaleX)
	valuesBuf.Write([]byte{0x00, 0x00, 0x80, 0x3f})

	// Without this, talent screen is blank
	updateMask.SetFieldMask(realmd.FieldMaskUnitBytes0)
	valuesBuf.WriteByte(byte(char.Race))
	valuesBuf.WriteByte(byte(char.Class))
	valuesBuf.WriteByte(byte(char.Gender))
	valuesBuf.WriteByte(byte(realmd.PowerTypeForClass(char.Class)))

	// Without this, character spawns in as a corpse
	updateMask.SetFieldMask(realmd.FieldMaskUnitHealth)
	valuesBuf.Write([]byte{100, 0, 0, 0})

	// Without this, UI doesn't show max health
	updateMask.SetFieldMask(realmd.FieldMaskUnitMaxHealth)
	valuesBuf.Write([]byte{100, 0, 0, 0})

	// Without this, character level appears as 0
	updateMask.SetFieldMask(realmd.FieldMaskUnitLevel)
	valuesBuf.Write([]byte{10, 0, 0, 0})

	// Without this, client segfaults
	updateMask.SetFieldMask(realmd.FieldMaskUnitFactionTemplate)
	valuesBuf.Write([]byte{byte(char.Race), 0, 0, 0})

	// Without this, client segfaults
	updateMask.SetFieldMask(realmd.FieldMaskUnitDisplayId)
	valuesBuf.Write([]byte{0x0C, 0x4D, 0x00, 0x00}) // human female

	// Without this, client segfaults
	updateMask.SetFieldMask(realmd.FieldMaskUnitNativeDisplayId)
	valuesBuf.Write([]byte{0x0C, 0x4D, 0x00, 0x00}) // human female

	mask := updateMask.Mask()
	inner.WriteByte(byte(len(mask)))
	binary.Write(&inner, binary.LittleEndian, mask)
	inner.Write(valuesBuf.Bytes())

	if err := client.SendPacketBytes(realmd.OpServerUpdateObject, inner.Bytes()); err != nil {
		return err
	}

	log.Println("sent spawn player")
	return nil
}
