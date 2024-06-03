package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type RealmService interface {
	Get(uint32) (*Realm, error)
	List() ([]*Realm, error)
	Create(*Realm) error
	Update(*Realm) (bool, error)
	Delete(uint32) (bool, error)
}

type DbRealmService struct {
	db *sqlx.DB
}

var _ RealmService = (*DbRealmService)(nil)

func NewDbRealmService(db *sqlx.DB) *DbRealmService {
	return &DbRealmService{db}
}

func (s *DbRealmService) Get(id uint32) (*Realm, error) {
	result := &Realm{}
	if err := s.db.Get(result, `SELECT * FROM realms WHERE id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbRealmService) List() ([]*Realm, error) {
	results := []*Realm{}
	if err := s.db.Select(&results, `SELECT * FROM realms`); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DbRealmService) Create(r *Realm) error {
	q := `
	INSERT INTO realms (name, type, host, region)
	VALUES (:name, :type, :host, :region)
	RETURNING id, created_at`
	result, err := s.db.NamedQuery(q, r)
	if err != nil {
		return err
	}
	result.Next()
	return result.StructScan(r)
}

func (s *DbRealmService) Update(r *Realm) (bool, error) {
	q := `UPDATE realms SET name=:name, type=:type, host=:host, region=:region WHERE id=:id`
	result, err := s.db.NamedExec(q, r)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}

func (s *DbRealmService) Delete(id uint32) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM realms WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}
