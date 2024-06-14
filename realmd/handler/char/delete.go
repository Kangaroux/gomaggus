package char

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/cmsg_char_delete.html
type deleteRequest struct {
	CharacterId uint64
}

// https://gtker.com/wow_messages/docs/smsg_char_delete.html#client-version-335
type deleteResponse struct {
	ResponseCode realmd.ResponseCode
}

func DeleteHandler(svc *realmd.Service, client *realmd.Client, data *realmd.ClientPacket) error {
	log.Println("start character delete")

	r := bytes.NewReader(data.Payload)
	p := deleteRequest{}
	if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(p.CharacterId))
	if err != nil {
		return err
	}

	deleted := false

	if char == nil {
		log.Println("client tried to delete non-existent character:", p.CharacterId)
	} else if char.AccountId != client.Account.Id {
		log.Println("client tried to delete character from another account:", p.CharacterId)
	} else if char.RealmId != client.Realm.Id {
		log.Println("client tried to delete character from another realm:", p.CharacterId)
	} else {
		deleted, err = svc.Chars.Delete(char.Id)

		if err != nil {
			return err
		}
	}

	resp := deleteResponse{}

	if deleted {
		resp.ResponseCode = realmd.RespCodeCharDeleteSuccess
	} else {
		resp.ResponseCode = realmd.RespCodeCharDeleteFailed
	}

	if err := client.SendPacket(realmd.OpServerCharCreate, &resp); err != nil {
		return err
	}

	log.Println("finished character delete")

	return nil
}
