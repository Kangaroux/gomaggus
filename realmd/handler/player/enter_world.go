package player

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/cmsg_player_login.html
type loginRequest struct {
	CharacterId uint64
}

func LoginHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	log.Println("start character login")

	r := bytes.NewReader(data[6:])
	p := loginRequest{}
	if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(p.CharacterId))
	if err != nil {
		return err
	}

	resp := bytes.Buffer{}
	ok := char != nil && char.AccountId == client.Account.Id && char.RealmId == client.Realm.Id

	if !ok {
		// https: gtker.com/wow_messages/docs/smsg_character_login_failed.html#client-version-335
		respHeader, err := client.BuildHeader(realmd.OpServerCharLoginFailed, 1)
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.WriteByte(byte(realmd.RespCodeCharLoginFailed))
	} else {
		// https://gtker.com/wow_messages/docs/smsg_login_verify_world.html
		inner := bytes.Buffer{}
		inner.Write([]byte{0, 0, 0, 0})                              // map (hardcoded as eastern kingdoms)
		binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // x
		binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // y
		binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // z
		binary.Write(&inner, binary.LittleEndian, float32(0))        // orientation

		respHeader, err := client.BuildHeader(realmd.OpServerCharLoginVerifyWorld, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.Write(inner.Bytes())
	}

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent verify world")

	if ok {
		// https://gtker.com/wow_messages/docs/smsg_tutorial_flags.html
		resp := bytes.Buffer{}
		respHeader, err := client.BuildHeader(realmd.OpServerTutorialFlags, 32)
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.Write(bytes.Repeat([]byte{255}, 32))

		if _, err := client.Conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent tutorial flags")

		// https://gtker.com/wow_messages/docs/smsg_feature_system_status.html#client-version-335
		inner := bytes.Buffer{}
		inner.WriteByte(2) // auto ignore?
		inner.WriteByte(0) // voip enabled

		resp = bytes.Buffer{}
		respHeader, err = client.BuildHeader(realmd.OpServerSystemFeatures, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.Write(inner.Bytes())

		if _, err := client.Conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent system features")

		// https://gtker.com/wow_messages/docs/smsg_bindpointupdate.html#client-version-335
		inner = bytes.Buffer{}
		binary.Write(&inner, binary.LittleEndian, float32(-8949.95)) // hearth x
		binary.Write(&inner, binary.LittleEndian, float32(-132.493)) // hearth y
		binary.Write(&inner, binary.LittleEndian, float32(83.5312))  // hearth z
		inner.Write([]byte{0, 0, 0, 0})                              // map: eastern kingdoms
		inner.Write([]byte{12, 0, 0, 0})                             // area: elwynn forest

		resp = bytes.Buffer{}
		respHeader, err = client.BuildHeader(realmd.OpServerHearthLocation, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.Write(inner.Bytes())

		if _, err := client.Conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent hearth location")

		// https://gtker.com/wow_messages/docs/smsg_trigger_cinematic.html
		// inner = bytes.Buffer{}
		// binary.Write(&inner, binary.LittleEndian, uint32(81)) // human

		// resp = bytes.Buffer{}
		// respHeader, err = client.BuildHeader(OpServerPlayCinematic, uint32(inner.Len()))
		// if err != nil {
		// 	return err
		// }
		// resp.Write(respHeader)
		// resp.Write(inner.Bytes())

		// if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		// 	return err
		// }

		// log.Println("sent play cinematic")

		// https://gtker.com/wow_messages/docs/smsg_update_object.html#client-version-335
		inner = bytes.Buffer{}
		inner.Write([]byte{1, 0, 0, 0}) // number of objects

		// nested object start
		inner.WriteByte(byte(realmd.UpdateTypeCreateObject2)) // update type: CREATE_OBJECT2
		inner.Write(internal.PackGuid(uint64(char.Id)))       // packed guid
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
		// field mask end

		// nested object end

		resp = bytes.Buffer{}
		respHeader, err = client.BuildHeader(realmd.OpServerUpdateObject, uint32(inner.Len()))
		if err != nil {
			return err
		}
		resp.Write(respHeader)
		resp.Write(inner.Bytes())

		fmt.Printf("%x\n", respHeader)
		fmt.Printf("%x\n", resp.Bytes())

		if _, err := client.Conn.Write(resp.Bytes()); err != nil {
			return err
		}

		log.Println("sent object update")
	}

	log.Println("finished character login")

	return nil
}
