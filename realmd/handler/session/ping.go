package session

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/cmsg_ping.html#client-version-19-client-version-110-client-version-111-client-version-112-client-version-2-client-version-3
type pingRequest struct {
	SequenceId    uint32
	RoundTripTime uint32 // zero if server hasn't responded?
}

func PingHandler(client *realmd.Client, data []byte) error {
	log.Println("starting ping")

	var err error
	r := bytes.NewReader(data[6:])
	p := pingRequest{}

	if err = binary.Read(r, binary.LittleEndian, &p.SequenceId); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.RoundTripTime); err != nil {
		return err
	}

	resp := bytes.Buffer{}
	respHeader, err := client.BuildHeader(realmd.OpServerPong, 4)
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	binary.Write(&resp, binary.LittleEndian, p.SequenceId)

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent pong")

	return nil
}
