package models

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CharacterListParams struct {
	AccountId uint32
	RealmId   uint32
}

type CharacterService interface {
	// Get returns the matching character, or nil if it doesn't exist.
	Get(uint32) (*Character, error)

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

func NewDbCharacterervice(db *sqlx.DB) *DbCharacterService {
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

	results := []*Character{}
	if err := s.db.Select(&results, q, args...); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *DbCharacterService) Create(r *Character) error {
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
		facial_hair,
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
		:facial_hair,
		:outfit_id
	) RETURNING id, created_at`
	result, err := s.db.NamedQuery(q, r)
	if err != nil {
		return err
	}
	result.Next()
	return result.StructScan(r)
}

func (s *DbCharacterService) Update(r *Character) (bool, error) {
	q := `
	UPDATE characters SET
		name=:name,
		race=:race,
		class=:class,
		gender=:gender,
		skin_color=:skin_color,
		face=:face,
		hair_style=:hair_style,
		hair_color=:hair_color,
		facial_hair=:facial_hair,
		outfit_id=:outfit_id
	WHERE
		id=:id`
	result, err := s.db.NamedExec(q, r)
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
