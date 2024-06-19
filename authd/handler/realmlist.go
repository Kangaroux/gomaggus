package handler

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/kangaroux/gomaggus/authd"
	"github.com/kangaroux/gomaggus/model"
	"github.com/mixcode/binarystruct"
)

type realmListRequest struct {
	_ authd.Opcode
	_ [4]byte // Padding
}

var realmListRequestSize = binary.Size(realmListRequest{})

// https://gtker.com/wow_messages/docs/cmd_realm_list_server.html#protocol-version-8
type realmListHeader struct {
	Opcode authd.Opcode // OpRealmList
	Size   uint16
}

type realmListBody struct {
	_         [4]byte // header padding
	NumRealms uint16
	Realms    []realm `binary:"[NumRealms]Any"`
	_         [2]byte // footer padding
}

type realm struct {
	Type          model.RealmType
	Locked        bool
	Flags         model.RealmFlag
	Name          string `binary:"zstring"`
	Host          string `binary:"zstring"`
	Population    float32
	NumCharacters uint8
	Region        model.RealmRegion
	Id            uint8
}

type RealmList struct {
	Client *authd.Client
	Realms model.RealmService
}

func (h *RealmList) Handle() error {
	if h.Client.State != authd.StateAuthenticated {
		return &ErrWrongState{
			Handler:  "RealmList",
			Expected: authd.StateAuthenticated,
			Actual:   h.Client.State,
		}
	}

	realmList, err := h.Realms.List()
	if err != nil {
		return err
	}

	respBody := realmListBody{
		NumRealms: uint16(len(realmList)),
		Realms:    make([]realm, len(realmList)),
	}

	for i, r := range realmList {
		respBody.Realms[i] = realm{
			Type:          r.Type,
			Locked:        false,
			Flags:         model.RealmFlagNone,
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

	respHeader := realmListHeader{
		Opcode: authd.OpcodeRealmList,
		Size:   uint16(len(bodyBytes)),
	}

	headerBytes, err := binarystruct.Marshal(&respHeader, binarystruct.LittleEndian)
	if err != nil {
		return err
	}

	respBuf := bytes.Buffer{}
	respBuf.Write(headerBytes)
	respBuf.Write(bodyBytes)

	if _, err := h.Client.Conn.Write(respBuf.Bytes()); err != nil {
		return err
	}

	log.Println("Replied to realm list")

	return nil
}

// Read verifies data is large enough, but does not use it. If data is too small, Read returns ErrPacketReadEOF.
func (h *RealmList) Read(data []byte) (int, error) {
	if len(data) < realmListRequestSize {
		return 0, ErrPacketReadEOF
	}
	return realmListRequestSize, nil
}
