package realmd

import "github.com/kangaroux/gomaggus/model"

type Service struct {
	Accounts       model.AccountService
	AccountStorage model.AccountStorageService
	Chars          model.CharacterService
	Realms         model.RealmService
	Sessions       model.SessionService
}
