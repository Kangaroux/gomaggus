package char

import (
	"encoding/binary"

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
		client.Log.Warn().Uint64("char", req.CharacterId).Msg("tried to delete non-existent character")
		return &realmd.ErrKickClient{Reason: "invalid char delete"}

	} else if char.AccountId != client.Account.Id {
		client.Log.Warn().Str("char", char.String()).Msg("tried to delete character from another account")
		return &realmd.ErrKickClient{Reason: "invalid char delete"}

	} else if char.RealmId != client.Realm.Id {
		client.Log.Warn().Str("char", char.String()).Msg("tried to delete character on another realm")
		return &realmd.ErrKickClient{Reason: "invalid char delete"}

	} else {
		deleted, err = svc.Chars.Delete(char.Id)

		if err != nil {
			return err
		}

		client.Log.Info().Str("char", char.String()).Msg("character deleted")
	}

	resp := deleteResponse{}

	if deleted {
		resp.ResponseCode = realmd.RespCodeCharDeleteSuccess
	} else {
		resp.ResponseCode = realmd.RespCodeCharDeleteFailed
	}

	return client.SendPacket(realmd.OpServerCharCreate, &resp)
}
