package realmd

import "github.com/kangaroux/gomaggus/model"

type Service struct {
	Accounts         model.AccountService
	AccountStorage   model.AccountStorageService
	CharacterStorage model.CharacterStorageService
	Characters       model.CharacterService
	Realms           model.RealmService
	Sessions         model.SessionService
}
