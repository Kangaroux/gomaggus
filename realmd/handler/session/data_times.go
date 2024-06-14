package session

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"

	"github.com/kangaroux/gomaggus/realmd"
)

func DataTimesHandler(client *realmd.Client) error {
	log.Println("starting account data times")

	inner := bytes.Buffer{}
	binary.Write(&inner, binary.LittleEndian, uint32(time.Now().Unix()))
	inner.WriteByte(1)                 // activated (bool)
	inner.Write([]byte{0, 0, 0, 0xFF}) // cache mask (all)
	// cache times
	for i := 0; i < 8; i++ {
		inner.Write([]byte{0, 0, 0, 0})
	}

	resp := bytes.Buffer{}
	respHeader, err := client.BuildHeader(realmd.OpServerAccountDataTimes, uint32(inner.Len()))
	if err != nil {
		return err
	}
	resp.Write(respHeader)
	resp.Write(inner.Bytes())

	if _, err := client.Conn.Write(resp.Bytes()); err != nil {
		return err
	}

	log.Println("sent account data times")

	return nil
}
