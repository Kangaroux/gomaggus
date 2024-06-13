package authd

import (
	"bytes"
	"log"

	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
type ServerRealm struct {
	Type          byte
	Locked        bool
	Flags         byte
	Name          string `binary:"zstring"`
	Host          string `binary:"zstring"`
	Population    float32
	NumCharacters byte
	Region        byte
	Id            byte
}

type ServerRealmListBody struct {
	_         [4]byte // header padding
	NumRealms uint16
	Realms    []ServerRealm `binary:"[NumRealms]Any"`
	_         [2]byte       // footer padding
}

type ServerRealmListHeader struct {
	Opcode byte
	Size   uint16
}

func handleRealmList(services *Services, c *Client) error {
	realmList, err := services.realms.List()
	if err != nil {
		return err
	}

	respBody := ServerRealmListBody{
		NumRealms: uint16(len(realmList)),
		Realms:    make([]ServerRealm, len(realmList)),
	}

	for i, r := range realmList {
		respBody.Realms[i] = ServerRealm{
			Type:          r.Type,
			Locked:        false,
			Flags:         RealmFlagNone,
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

	respHeader := ServerRealmListHeader{
		Opcode: OpRealmList,
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
