package authd

import (
	"bytes"
	"encoding/binary"
	"log"
)

func handleRealmList(services *Services, c *Client) error {
	realmList, err := services.realms.List()
	if err != nil {
		return err
	}

	resp := &bytes.Buffer{}
	resp.WriteByte(OP_REALM_LIST)

	inner := &bytes.Buffer{}
	inner.Write([]byte{0, 0, 0, 0}) // header padding
	binary.Write(inner, binary.LittleEndian, uint16(len(realmList)))
	for _, r := range realmList {
		inner.WriteByte(byte(r.Type))
		inner.WriteByte(0)                    // locked
		inner.WriteByte(byte(REALMFLAG_NONE)) // TODO?
		inner.WriteString(r.Name)
		inner.WriteByte(0) // name is NUL-terminated
		inner.WriteString(r.Host)
		inner.WriteByte(0)                                   // host is NUL-terminated
		binary.Write(inner, binary.LittleEndian, float32(0)) // TODO: population
		inner.WriteByte(byte(0))                             // TODO: number of chars on realm
		inner.WriteByte(byte(r.Region))
		inner.WriteByte(byte(r.Id))
	}
	inner.Write([]byte{0, 0}) // footer padding

	// Write size of realm list payload
	binary.Write(resp, binary.LittleEndian, uint16(inner.Len()))
	// Concat to main payload
	inner.WriteTo(resp)

	if _, err := c.conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to realm list")

	return nil
}
