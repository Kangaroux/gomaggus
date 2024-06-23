package char

import (
	"github.com/kangaroux/gomaggus/model"
	"github.com/kangaroux/gomaggus/realmd"
)

// https://gtker.com/wow_messages/docs/charactergear.html
type gearDisplay struct {
	DisplayId uint32
	// https://gtker.com/wow_messages/docs/inventorytype.html
	Slot        uint8
	Enchantment uint32
}

// https://gtker.com/wow_messages/docs/character.html#client-version-335
type character struct {
	Guid                 realmd.Guid
	Name                 string `binary:"zstring"`
	Race                 model.Race
	Class                model.Class
	Gender               model.Gender
	Skin                 uint8
	Face                 uint8
	HairStyle            uint8
	HairColor            uint8
	ExtraCosmetic        uint8
	Level                uint8
	Area                 uint32
	Map                  uint32
	Position             realmd.Vector3
	GuildId              uint32
	Flags                uint32
	RecustomizationFlags uint32
	FirstLogin           bool
	PetDisplayId         uint32
	PetLevel             uint32
	PetFamily            uint32
	GearDisplay          [23]gearDisplay
}

// https://gtker.com/wow_messages/docs/smsg_char_enum.html#client-version-335
type listResponse struct {
	Count      uint8
	Characters []character `binary:"[Count]Any"`
}

func ListHandler(svc *realmd.Service, client *realmd.Client) error {
	accountChars, err := svc.Characters.List(&model.CharacterListParams{
		AccountId: client.Account.Id,
		RealmId:   client.Realm.Id,
		Sort:      model.OldestToNewest,
	})
	if err != nil {
		return err
	}

	resp := listResponse{
		Count:      uint8(len(accountChars)),
		Characters: make([]character, len(accountChars)),
	}

	for i, accountChar := range accountChars {
		char := character{
			Guid:                 realmd.Guid(accountChar.Id),
			Name:                 accountChar.Name,
			Race:                 accountChar.Race,
			Class:                accountChar.Class,
			Gender:               accountChar.Gender,
			Skin:                 accountChar.SkinColor,
			Face:                 accountChar.Face,
			HairStyle:            accountChar.HairStyle,
			HairColor:            accountChar.HairStyle,
			ExtraCosmetic:        accountChar.ExtraCosmetic,
			Level:                1,
			Area:                 0xC,              // Elwynn forest
			Map:                  0x0,              // Eastern kingdoms
			Position:             realmd.Vector3{}, // Position doesn't matter for char list
			GuildId:              0,
			Flags:                0, // ??
			RecustomizationFlags: 0, // ??
			FirstLogin:           !accountChar.LastLogin.Valid,
			PetDisplayId:         0,
			PetLevel:             0,
			PetFamily:            0,
		}

		for j := 0; j < 23; j++ {
			char.GearDisplay[j] = gearDisplay{
				DisplayId:   0,
				Slot:        uint8(j) + 1,
				Enchantment: 0,
			}
		}

		resp.Characters[i] = char
	}

	return client.SendPacket(realmd.OpServerCharList, &resp)
}
