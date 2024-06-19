package mock

import "github.com/kangaroux/gomaggus/model"

type SessionService struct {
	OnGet            func(uint32) (*model.Session, error)
	OnCreate         func(*model.Session) error
	OnUpdate         func(*model.Session) (bool, error)
	OnDelete         func(uint32) (bool, error)
	OnUpdateOrCreate func(*model.Session) (bool, error)
}

var _ model.SessionService = (*SessionService)(nil)

func (s *SessionService) Get(id uint32) (*model.Session, error) {
	if s.OnGet == nil {
		return nil, nil
	}
	return s.OnGet(id)
}

func (s *SessionService) Create(session *model.Session) error {
	if s.OnCreate == nil {
		return nil
	}
	return s.OnCreate(session)
}

func (s *SessionService) Update(session *model.Session) (bool, error) {
	if s.OnUpdate == nil {
		return true, nil
	}
	return s.OnUpdate(session)
}

func (s *SessionService) Delete(id uint32) (bool, error) {
	if s.OnDelete == nil {
		return true, nil
	}
	return s.OnDelete(id)
}

func (s *SessionService) UpdateOrCreate(session *model.Session) (bool, error) {
	if s.OnUpdateOrCreate == nil {
		return false, nil
	}
	return s.OnUpdateOrCreate(session)
}
