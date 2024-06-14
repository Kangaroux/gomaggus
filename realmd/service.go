package realmd

import "github.com/kangaroux/gomaggus/model"

type Service struct {
	Accounts model.AccountService
	Chars    model.CharacterService
	Realms   model.RealmService
	Sessions model.SessionService
}
