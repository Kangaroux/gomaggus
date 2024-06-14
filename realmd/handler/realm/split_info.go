package realm

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
)

type RealmSplitState uint32

const (
	SplitNormal    = 0
	SplitConfirmed = 1
	SplitPending   = 2
)

// https://gtker.com/wow_messages/docs/cmsg_realm_split.html
type splitRequest struct {
	RealmId uint32
}

func SplitInfoHandler(client *realmd.Client, data []byte) error {
	log.Println("starting realm split")

	r := bytes.NewReader(data[6:])
	p := splitRequest{}
	binary.Read(r, binary.LittleEndian, &p.RealmId)

	// https://gtker.com/wow_messages/docs/smsg_realm_split.html
	inner := bytes.Buffer{}
	binary.Write(&inner, binary.LittleEndian, p.RealmId)
	binary.Write(&inner, binary.LittleEndian, uint32(SplitNormal))
	inner.WriteString("01/01/01\x00") // send a bogus date (NUL-terminated)

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(realmd.OpServerRealmSplit, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(client.Crypto.Encrypt(respHeader))
	resp.Write(inner.Bytes())

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent realm split")

	return nil
}
