package char

import (
	"log"

	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
	"github.com/mixcode/binarystruct"
)

// https://gtker.com/wow_messages/docs/cmsg_char_create.html#client-version-32-client-version-33
type createRequest struct {
	Name          string `binary:"zstring"`
	Race          model.Race
	Class         model.Class
	Gender        model.Gender
	SkinColor     byte
	Face          byte
	HairStyle     byte
	HairColor     byte
	ExtraCosmetic byte
	OutfitId      byte
}

// https://gtker.com/wow_messages/docs/smsg_char_create.html#client-version-335
type createResponse struct {
	ResponseCode realmd.ResponseCode
}

func CreateHandler(svc *realmd.Service, client *realmd.Client, data []byte) error {
	// TODO: check if account is full
	// accountChars, err := s.charsDb.List(&model.CharacterListParams{
	// 	AccountId: c.account.Id,
	// 	RealmId:   c.realm.Id,
	// })
	// if err != nil {
	// 	return err
	// }

	req := createRequest{}
	if _, err := binarystruct.Unmarshal(data, binarystruct.LittleEndian, &req); err != nil {
		return err
	}

	log.Println("client wants to create character", req.Name)

	existing, err := svc.Chars.GetName(req.Name, client.Realm.Id)
	if err != nil {
		return err
	}

	resp := createResponse{}

	if existing != nil {
		resp.ResponseCode = realmd.RespCodeCharCreateNameInUse
	} else {
		char := &model.Character{
			Name:          req.Name,
			AccountId:     client.Account.Id,
			RealmId:       client.Realm.Id,
			Race:          req.Race,
			Class:         req.Class,
			Gender:        req.Gender,
			SkinColor:     req.SkinColor,
			Face:          req.Face,
			HairStyle:     req.HairStyle,
			HairColor:     req.HairColor,
			ExtraCosmetic: req.ExtraCosmetic,
			OutfitId:      req.OutfitId,
		}
		if err := svc.Chars.Create(char); err != nil {
			return err
		}

		log.Println("Created new", char)
		resp.ResponseCode = realmd.RespCodeCharCreateSuccess
	}

	return client.SendPacket(realmd.OpServerCharCreate, &resp)
}
