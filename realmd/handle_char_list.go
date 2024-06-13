package realmd

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/models"
)

func handleCharList(services *Services, client *Client) error {
	log.Println("starting character list")

	accountChars, err := services.chars.List(&models.CharacterListParams{
		AccountId: client.account.Id,
		RealmId:   client.realm.Id,
	})
	if err != nil {
		return err
	}

	// https://gtker.com/wow_messages/docs/smsg_char_enum.html#client-version-335
	inner := bytes.Buffer{}
	inner.WriteByte(byte(len(accountChars)))

	for _, char := range accountChars {
		binary.Write(&inner, binary.LittleEndian, uint64(char.Id))
		inner.WriteString(char.Name)
		inner.WriteByte(0) // NUL-terminated
		inner.WriteByte(byte(char.Race))
		inner.WriteByte(byte(char.Class))
		inner.WriteByte(byte(char.Gender))
		inner.WriteByte(char.SkinColor)
		inner.WriteByte(char.Face)
		inner.WriteByte(char.HairStyle)
		inner.WriteByte(char.HairColor)
		inner.WriteByte(char.FacialHair)
		inner.WriteByte(1)                                    // level
		inner.Write([]byte{12, 0, 0, 0})                      // area (hardcoded as elwynn forest)
		inner.Write([]byte{0, 0, 0, 0})                       // map (hardcoded as eastern kingdoms)
		binary.Write(&inner, binary.LittleEndian, float32(0)) // x
		binary.Write(&inner, binary.LittleEndian, float32(0)) // y
		binary.Write(&inner, binary.LittleEndian, float32(0)) // z
		inner.Write([]byte{0, 0, 0, 0})                       // guild id
		inner.Write([]byte{0, 0, 0, 0})                       // flags
		inner.Write([]byte{0, 0, 0, 0})                       // recustomization_flags (?)

		if !char.LastLogin.Valid {
			inner.WriteByte(1) // first login, show tutorial
		} else {
			inner.WriteByte(0) // not first login
		}

		inner.Write([]byte{0, 0, 0, 0}) // pet display id
		inner.Write([]byte{0, 0, 0, 0}) // pet level
		inner.Write([]byte{0, 0, 0, 0}) // pet family

		// equipment (from head to holdable)
		// https://gtker.com/wow_messages/docs/inventorytype.html
		for i := 1; i <= 23; i++ {
			inner.Write([]byte{0, 0, 0, 0}) // equipment display id
			inner.WriteByte(byte(i))        // slot
			inner.Write([]byte{0, 0, 0, 0}) // enchantment
		}
	}

	resp := bytes.Buffer{}
	respHeader, err := makeServerHeader(OpServerCharEnum, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))
	resp.Write(inner.Bytes())

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent character list")

	return nil
}
