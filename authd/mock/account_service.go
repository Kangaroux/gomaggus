package mock

import "github.com/kangaroux/gomaggus/model"

type AccountService struct {
	OnGet    func(*model.AccountGetParams) (*model.Account, error)
	OnList   func() ([]*model.Account, error)
	OnCreate func(*model.Account) error
	OnUpdate func(*model.Account) (bool, error)
	OnDelete func(uint32) (bool, error)
}

var _ model.AccountService = (*AccountService)(nil)

func (s *AccountService) Get(params *model.AccountGetParams) (*model.Account, error) {
	if s.OnGet == nil {
		return nil, nil
	}
	return s.OnGet(params)
}

func (s *AccountService) List() ([]*model.Account, error) {
	if s.OnList == nil {
		return nil, nil
	}
	return s.OnList()
}

func (s *AccountService) Create(acct *model.Account) error {
	if s.OnCreate == nil {
		return nil
	}
	return s.OnCreate(acct)
}

func (s *AccountService) Update(acct *model.Account) (bool, error) {
	if s.OnUpdate == nil {
		return true, nil
	}
	return s.OnUpdate(acct)
}

func (s *AccountService) Delete(id uint32) (bool, error) {
	if s.OnDelete == nil {
		return true, nil
	}
	return s.OnDelete(id)
}
