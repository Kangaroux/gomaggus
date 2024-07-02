package player

import (
	"encoding/binary"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

type nameRequestError uint8

const (
	noError  nameRequestError = 0
	notFound nameRequestError = 1
)

type nameRequest struct {
	Guid realmd.Guid
}

type nameResponseNotFound struct {
	Guid      realmd.PackedGuid
	ErrorCode nameRequestError
}

type nameResponseOK struct {
	Guid             realmd.PackedGuid
	ErrorCode        nameRequestError
	CharacterName    string `binary:"zstring"`
	RealmName        string `binary:"zstring"`
	Race             model.Race
	Gender           model.Gender
	Class            model.Class
	HasDeclinedNames bool // Absolutely no idea what this is

	// Should be included if HasDeclinedNames == true
	// DeclinedNames [5]string `binary:"[5]zstring"`
}

func NameHandler(svc realmd.Service, client *realmd.Client, data []byte) error {
	req := nameRequest{}
	if _, err := binarystruct.Unmarshal(data, binary.LittleEndian, &req); err != nil {
		return err
	}

	char, err := svc.Characters.Get(uint32(req.Guid))
	if err != nil {
		return err
	}

	if char == nil {
		resp := nameResponseNotFound{
			Guid:      realmd.PackGuid(uint64(req.Guid)),
			ErrorCode: notFound,
		}
		return client.SendPacket(realmd.OpServerGetPlayerNameResponse, &resp)
	}

	var realmName string

	if char.RealmId != client.Realm.Id {
		realm, err := svc.Realms.Get(char.RealmId)
		if err != nil {
			return err
		}

		realmName = realm.Name
	} else {
		realmName = client.Realm.Name
	}

	resp := nameResponseOK{
		Guid:             realmd.PackGuid(uint64(req.Guid)),
		ErrorCode:        noError,
		CharacterName:    char.Name,
		RealmName:        realmName,
		Race:             char.Race,
		Gender:           char.Gender,
		Class:            char.Class,
		HasDeclinedNames: false,
	}
	return client.SendPacket(realmd.OpServerGetPlayerNameResponse, &resp)
}
