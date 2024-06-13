package model

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
	// Get returns an account, or nil if it doesn't exist. At least one param must be specified.
	// Params are combined using OR.
	Get(*AccountGetParams) (*Account, error)

	// List returns a list of all accounts.
	List() ([]*Account, error)

	// Create creates a new account and sets the Id and CreatedAt fields.
	Create(*Account) error

	// Update tries to update an existing account and returns if it was updated.
	Update(*Account) (bool, error)

	// Delete tries to delete an existing account by id and returns if it was deleted.
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
	q := `SELECT * FROM accounts WHERE `
	cond := []string{}
	args := []interface{}{}

	if params.Id > 0 {
		cond = append(cond, `id = ?`)
		args = append(args, params.Id)
	}
	if len(params.Email) > 0 {
		cond = append(cond, `lower(email) = ?`)
		args = append(args, strings.ToLower(params.Email))
	}
	if len(params.Username) > 0 {
		cond = append(cond, `upper(username) = ?`)
		args = append(args, strings.ToUpper(params.Username))
	}
	if len(cond) == 0 {
		return nil, ErrEmptyGetParams
	}

	q += strings.Join(cond, " OR ")
	q = s.db.Rebind(q)

	result := &Account{}
	if err := s.db.Get(result, q, args...); err != nil {
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
	INSERT INTO accounts (username, email, srp_verifier, srp_salt)
	VALUES (:username, :email, :srp_verifier, :srp_salt)
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
