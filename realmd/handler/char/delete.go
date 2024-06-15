package char

import (
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
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
	req := deleteRequest{}
	if _, err := binarystruct.Unmarshal(data.Payload, binary.LittleEndian, &req); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(req.CharacterId))
	if err != nil {
		return err
	}

	deleted := false

	if char == nil {
		log.Println("client tried to delete non-existent character:", req.CharacterId)
	} else if char.AccountId != client.Account.Id {
		log.Println("client tried to delete character from another account:", req.CharacterId)
	} else if char.RealmId != client.Realm.Id {
		log.Println("client tried to delete character from another realm:", req.CharacterId)
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
