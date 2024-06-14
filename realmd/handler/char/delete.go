package char

import (
	"bytes"
	"encoding/binary"
	"log"
)

func handleCharDelete(services *Services, client *Client, data []byte) error {
	log.Println("start character delete")

	r := bytes.NewReader(data[6:])
	p := CharDeletePacket{}
	if err := binary.Read(r, binary.LittleEndian, &p.CharacterId); err != nil {
		return err
	}

	char, err := services.chars.Get(uint32(p.CharacterId))
	if err != nil {
		return err
	}

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(OpServerCharDelete, 1)
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))

	if char == nil || char.AccountId != client.account.Id || char.RealmId != client.realm.Id {
		resp.WriteByte(byte(RespCodeCharDeleteFailed))
	} else {
		if _, err := services.chars.Delete(char.Id); err != nil {
			return err
		}
		resp.WriteByte(byte(RespCodeCharDeleteSuccess))
	}

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("finished character create")

	return nil
}
