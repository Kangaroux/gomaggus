package session

import (
	"bytes"
	"encoding/binary"
	"log"
)

func handlePing(client *Client, data []byte) error {
	log.Println("starting ping")

	var err error
	r := bytes.NewReader(data[6:])
	p := PingPacket{}

	if err = binary.Read(r, binary.LittleEndian, &p.SequenceId); err != nil {
		return err
	}
	if err = binary.Read(r, binary.LittleEndian, &p.RoundTripTime); err != nil {
		return err
	}

	resp := bytes.Buffer{}
	respHeader, err := realmd.BuildHeader(OpServerPong, 4)
	if err != nil {
		return err
	}
	resp.Write(client.crypto.Encrypt(respHeader))
	binary.Write(&resp, binary.LittleEndian, p.SequenceId)

	if _, err := client.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent pong")

	return nil
}
