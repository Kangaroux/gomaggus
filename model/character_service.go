package model

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CharacterSort uint8

const (
	NoSort CharacterSort = iota
	OldestToNewest
	Alphabetically
	LastLogin
)

type CharacterListParams struct {
	AccountId uint32
	RealmId   uint32
	Sort      CharacterSort
}

type CharacterService interface {
	// Get returns a character by id, or nil if it doesn't exist.
	Get(uint32) (*Character, error)

	// GetName returns a character by its name and realm id, or nil if it doesn't exist.
	GetName(name string, realmId uint32) (*Character, error)

	// List returns a list of all characters matching the search query. Any number of params can be
	// specified. Params are combined using AND.
	List(*CharacterListParams) ([]*Character, error)

	// Create creates a new character and sets the Id and CreatedAt fields.
	Create(*Character) error

	// Update tries to update an existing character and returns if it was updated.
	Update(*Character) (bool, error)

	// Delete tries to delete an existing character by id and returns if it was deleted.
	Delete(uint32) (bool, error)
}

type DbCharacterService struct {
	db *sqlx.DB
}

var _ CharacterService = (*DbCharacterService)(nil)

func NewDbCharacterService(db *sqlx.DB) *DbCharacterService {
	return &DbCharacterService{db}
}

func (s *DbCharacterService) Get(id uint32) (*Character, error) {
	result := &Character{}
	if err := s.db.Get(result, `SELECT * FROM characters WHERE id = $1`, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbCharacterService) GetName(name string, realmId uint32) (*Character, error) {
	q := `SELECT * FROM characters WHERE lower(name) = $1 AND realm_id = $2`
	result := &Character{}
	if err := s.db.Get(result, q, strings.ToLower(name), realmId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbCharacterService) List(params *CharacterListParams) ([]*Character, error) {
	q := `SELECT * FROM characters`
	cond := []string{}
	args := []interface{}{}

	if params.AccountId > 0 {
		cond = append(cond, `account_id = ?`)
		args = append(args, params.AccountId)
	}
	if params.RealmId > 0 {
		cond = append(cond, `realm_id = ?`)
		args = append(args, params.RealmId)
	}
	if len(cond) > 0 {
		q += " WHERE " + strings.Join(cond, " AND ")
		q = s.db.Rebind(q)
	}

	switch params.Sort {
	case OldestToNewest:
		q += " ORDER BY created_at ASC"
	case Alphabetically:
		q += " ORDER BY name ASC"
	case LastLogin:
		q += " ORDER BY last_login DESC NULLS LAST"
	}

	results := []*Character{}
	if err := s.db.Select(&results, q, args...); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DbCharacterService) Create(c *Character) error {
	q := `
	INSERT INTO characters (
		name,
		account_id,
		realm_id,
		race,
		class,
		gender,
		skin_color,
		face,
		hair_style,
		hair_color,
		extra_cosmetic,
		outfit_id
	) VALUES (
		:name,
		:account_id,
		:realm_id,
		:race,
		:class,
		:gender,
		:skin_color,
		:face,
		:hair_style,
		:hair_color,
		:extra_cosmetic,
		:outfit_id
	) RETURNING id, created_at`
	result, err := s.db.NamedQuery(q, c)
	if err != nil {
		return err
	}
	result.Next()
	return result.StructScan(c)
}

func (s *DbCharacterService) Update(c *Character) (bool, error) {
	q := `
	UPDATE characters SET
		last_login=:last_login,
		name=:name,
		race=:race,
		class=:class,
		gender=:gender,
		skin_color=:skin_color,
		face=:face,
		hair_style=:hair_style,
		hair_color=:hair_color,
		extra_cosmetic=:extra_cosmetic,
		outfit_id=:outfit_id
	WHERE
		id=:id`
	result, err := s.db.NamedExec(q, c)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}

func (s *DbCharacterService) Delete(id uint32) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM characters WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}
