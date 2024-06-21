package char

import (
	"encoding/binary"
	golog "log"

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

func DeleteHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	req := deleteRequest{}
	if _, err := binarystruct.Unmarshal(data, binary.LittleEndian, &req); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(req.CharacterId))
	if err != nil {
		return err
	}

	deleted := false

	if char == nil {
		golog.Println("client tried to delete non-existent character:", req.CharacterId)
	} else if char.AccountId != client.Account.Id {
		golog.Println("client tried to delete character from another account:", req.CharacterId)
	} else if char.RealmId != client.Realm.Id {
		golog.Println("client tried to delete character from another realm:", req.CharacterId)
	} else {
		deleted, err = svc.Chars.Delete(char.Id)

		if err != nil {
			return err
		}

		golog.Println("Deleted", char)
	}

	resp := deleteResponse{}

	if deleted {
		resp.ResponseCode = realmd.RespCodeCharDeleteSuccess
	} else {
		resp.ResponseCode = realmd.RespCodeCharDeleteFailed
	}

	return client.SendPacket(realmd.OpServerCharCreate, &resp)
}
