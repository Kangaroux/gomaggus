package authd

import "github.com/kangaroux/gomaggus/model"

type Service struct {
	Accounts model.AccountService
	Realms   model.RealmService
	Sessions model.SessionService
}
