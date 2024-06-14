package realm

import (
	"bytes"
	"encoding/binary"
	"log"
)

func handleRealmSplit(client *Client, data []byte) error {
	log.Println("starting realm split")

	r := bytes.NewReader(data[6:])
	p := RealmSplitPacket{}
	binary.Read(r, binary.LittleEndian, &p.RealmId)

	// https://gtker.com/wow_messages/docs/smsg_realm_split.html
	inner := bytes.Buffer{}
	binary.Write(&inner, binary.LittleEndian, p.RealmId)
	inner.Write([]byte{0, 0, 0, 0})   // split state, 0 = normal
	inner.WriteString("01/01/01\x00") // send a bogus date (NUL-terminated)

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(OpServerRealmSplit, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))
	resp.Write(inner.Bytes())

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent realm split")

	return nil
}
