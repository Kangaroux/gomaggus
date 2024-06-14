package char

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"

	"github.com/kangaroux/gomaggus/internal"
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/cmsg_char_create.html#client-version-32-client-version-33
type createRequest struct {
	// Name string
	Race       model.Race
	Class      model.Class
	Gender     model.Gender
	SkinColor  byte
	Face       byte
	HairStyle  byte
	HairColor  byte
	FacialHair byte
	OutfitId   byte
}

// https://gtker.com/wow_messages/docs/smsg_char_create.html#client-version-335
type createResponse struct {
	ResponseCode realmd.ResponseCode
}

func CreateHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	log.Println("starting character create")

	// TODO: check if account is full
	// accountChars, err := s.charsDb.List(&model.CharacterListParams{
	// 	AccountId: c.account.Id,
	// 	RealmId:   c.realm.Id,
	// })
	// if err != nil {
	// 	return err
	// }

	p := createRequest{}
	r := bytes.NewReader(data[6:])
	charName, err := internal.ReadCString(r)
	if err != nil {
		return err
	}
	charName = strings.TrimSpace(charName)

	if err := binary.Read(r, binary.BigEndian, &p); err != nil {
		return err
	}

	log.Println("client wants to create character", charName)

	existing, err := svc.Chars.GetName(charName, client.Realm.Id)
	if err != nil {
		return err
	}

	if existing == nil {
		char := &model.Character{
			Name:       charName,
			AccountId:  client.Account.Id,
			RealmId:    client.Realm.Id,
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
		if err := svc.Chars.Create(char); err != nil {
			return err
		}
		log.Println("created char with id", char.Id)
	}

	resp := createResponse{}

	if existing == nil {
		resp.ResponseCode = realmd.RespCodeCharCreateSuccess
	} else {
		resp.ResponseCode = realmd.RespCodeCharCreateNameInUse
	}

	if err := client.SendPacket(realmd.OpServerCharCreate, &resp); err != nil {
		return err
	}

	log.Println("finished character create")

	return nil
}
