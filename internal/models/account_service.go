package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AccountGetParams struct {
	Id       uint32
	Email    string
	Username string
}

var ErrEmptyGetParams = errors.New("at least one get param must be set")

type AccountService interface {
	Get(*AccountGetParams) (*Account, error)
	List() ([]*Account, error)
	Create(*Account) error
	Update(*Account) (bool, error)
	Delete(uint32) (bool, error)
}

type DbAccountService struct {
	db *sqlx.DB
}

var _ AccountService = (*DbAccountService)(nil)

func NewDbAccountService(db *sqlx.DB) *DbAccountService {
	return &DbAccountService{db}
}

func (s *DbAccountService) Get(params *AccountGetParams) (*Account, error) {
	var arg interface{}
	q := `SELECT * FROM accounts WHERE `

	if params.Id > 0 {
		q += `id = $1`
		arg = params.Id
	} else if len(params.Email) > 0 {
		q += `lower(email) = $1`
		arg = strings.ToLower(params.Email)
	} else if len(params.Username) > 0 {
		q += `upper(username) = $1`
		arg = strings.ToUpper(params.Username)
	} else {
		return nil, ErrEmptyGetParams
	}

	result := &Account{}
	if err := s.db.Get(result, q, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbAccountService) List() ([]*Account, error) {
	results := []*Account{}
	if err := s.db.Select(&results, `SELECT * FROM accounts`); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DbAccountService) Create(a *Account) error {
	q := `
	INSERT INTO accounts (username, email, srp_verifier, srp_salt, realm_id)
	VALUES (:username, :email, :srp_verifier, :srp_salt, :realm_id)
	RETURNING id, created_at`
	result, err := s.db.NamedQuery(q, a)
	if err != nil {
		return err
	}
	result.Next()
	return result.StructScan(a)
}

func (s *DbAccountService) Update(a *Account) (bool, error) {
	q := `
	UPDATE accounts SET
	username=:username, email=:email, srp_verifier=:srp_verifier, srp_salt=:srp_salt, last_login=:last_login
	WHERE id=:id`
	result, err := s.db.NamedExec(q, a)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}

func (s *DbAccountService) Delete(id uint32) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM accounts WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}
