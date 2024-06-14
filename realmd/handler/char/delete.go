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

func DeleteHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	log.Println("start character delete")

	r := bytes.NewReader(data[6:])
	p := deleteRequest{}
	if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
		return err
	}

	char, err := svc.Chars.Get(uint32(p.CharacterId))
	if err != nil {
		return err
	}

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(realmd.OpServerCharDelete, 1)
	if err != nil {
		return err
	}
	resp.Write(client.Crypto.Encrypt(respHeader))

	if char == nil || char.AccountId != client.Account.Id || char.RealmId != client.Realm.Id {
		resp.WriteByte(byte(realmd.RespCodeCharDeleteFailed))
	} else {
		if _, err := svc.Chars.Delete(char.Id); err != nil {
			return err
		}
		resp.WriteByte(byte(realmd.RespCodeCharDeleteSuccess))
	}

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("finished character create")

	return nil
}
