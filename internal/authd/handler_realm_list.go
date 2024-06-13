package authd

import (
	"bytes"
	"log"

	"github.com/kangaroux/gomaggus/internal/authd/packets"
	"github.com/mixcode/binarystruct"
)

func handleRealmList(services *Services, c *Client) error {
	realmList, err := services.realms.List()
	if err != nil {
		return err
	}

	respBody := packets.ServerRealmListBody{
		NumRealms: uint16(len(realmList)),
		Realms:    make([]packets.ServerRealm, len(realmList)),
	}

	for i, r := range realmList {
		respBody.Realms[i] = packets.ServerRealm{
			Type:          r.Type,
			Locked:        false,
			Flags:         REALMFLAG_NONE,
			Name:          r.Name,
			Host:          r.Host,
			Population:    0, // TODO
			NumCharacters: 0, // TODO
			Region:        r.Region,
			Id:            byte(r.Id),
		}
	}

	bodyBytes, err := binarystruct.Marshal(&respBody, binarystruct.LittleEndian)
	if err != nil {
		return err
	}

	respHeader := packets.ServerRealmListHeader{
		Opcode: OP_REALM_LIST,
		Size:   uint16(len(bodyBytes)),
	}

	headerBytes, err := binarystruct.Marshal(&respHeader, binarystruct.LittleEndian)
	if err != nil {
		return err
	}

	respBuf := bytes.Buffer{}
	respBuf.Write(headerBytes)
	respBuf.Write(bodyBytes)

	if _, err := c.conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to realm list")

	return nil
}
