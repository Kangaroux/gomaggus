package char

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"

	"github.com/kangaroux/gomaggus/model"
)

func handleCharCreate(services *Services, client *Client, data []byte) error {
	log.Println("starting character create")

	// TODO: check if account is full
	// accountChars, err := s.charsDb.List(&model.CharacterListParams{
	// 	AccountId: c.account.Id,
	// 	RealmId:   c.realm.Id,
	// })
	// if err != nil {
	// 	return err
	// }

	p := CharCreatePacket{}
	r := bytes.NewReader(data[6:])
	charName, err := readCString(r)
	if err != nil {
		return err
	}
	charName = strings.TrimSpace(charName)

	if err := binary.Read(r, binary.BigEndian, &p); err != nil {
		return err
	}

	log.Println("client wants to create character", charName)

	existing, err := services.chars.GetName(charName, client.realm.Id)
	if err != nil {
		return err
	}

	if existing == nil {
		char := &model.Character{
			Name:       charName,
			AccountId:  client.account.Id,
			RealmId:    client.realm.Id,
			Race:       p.Race,   // TODO
			Class:      p.Class,  // TODO
			Gender:     p.Gender, // TODO
			SkinColor:  p.SkinColor,
			Face:       p.Face,
			HairStyle:  p.HairStyle,
			HairColor:  p.HairColor,
			FacialHair: p.FacialHair,
			OutfitId:   p.OutfitId,
		}
		if err := services.chars.Create(char); err != nil {
			return err
		}
		log.Println("created char with id", char.Id)
	}

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(OpServerCharCreate, 1)
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))

	if existing != nil {
		resp.WriteByte(byte(RespCodeCharCreateNameInUse))
	} else {
		resp.WriteByte(byte(RespCodeCharCreateSuccess))
	}

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("finished character create")

	return nil
}
