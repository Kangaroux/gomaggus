package models

import (
	"github.com/jmoiron/sqlx"
)

type RealmService interface {
	List() ([]*Realm, error)
	Create(*Realm) error
	Update(*Realm) error
	Delete(uint32) error
}

type DbRealmService struct {
	db *sqlx.DB
}

var _ RealmService = (*DbRealmService)(nil)

func NewDbRealmService(db *sqlx.DB) *DbRealmService {
	return &DbRealmService{db}
}

func (s *DbRealmService) List() ([]*Realm, error) {
	results := []*Realm{}
	if err := s.db.Select(&results, `SELECT * FROM realms`); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DbRealmService) Create(r *Realm) error {
	result, err := s.db.Exec(
		`INSERT INTO realms (name, type, host, region)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		r.Name, r.Type, r.Host, r.Region,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	r.Id = uint32(id)
	return nil
}

func (s *DbRealmService) Update(r *Realm) error {
	_, err := s.db.Exec(
		`UPDATE realms SET name=$1 type=$2 host=$3 region=$4 WHERE id=$5`,
		r.Name, r.Type, r.Host, r.Region, r.Id,
	)
	return err
}

func (s *DbRealmService) Delete(id uint32) error {
	_, err := s.db.Exec(`DELETE FROM realms WHERE id=$1`, id)
	return err
}
