package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AccountStorageService interface {
	// Get returns the account storage with the given type. If it doesn't exist, the returned storag is nil.
	Get(uint32, AccountStorageType) (*AccountStorage, error)

	// List returns all storage belonging to accountID whose type matches mask. List may return fewer
	// results than specified in the mask. The results are not ordered.
	List(uint32, uint8) ([]*AccountStorage, error)

	// UpdateOrCreate tries to update the storage if it exists, otherwise it's created. UpdateOrCreate
	// reports whether the storage was created. If the error is nil and created is false, then it
	// was updated.
	UpdateOrCreate(*AccountStorage) (bool, error)
}

type DbAccountStorageService struct {
	db *sqlx.DB
}

var _ AccountStorageService = (*DbAccountStorageService)(nil)

func NewDbAccountStorageService(db *sqlx.DB) AccountStorageService {
	return &DbAccountStorageService{db}
}

func (s *DbAccountStorageService) Get(accountID uint32, storageType AccountStorageType) (*AccountStorage, error) {
	result := &AccountStorage{}
	q := `SELECT * FROM account_storage WHERE account_id = $1 AND type = $2`
	if err := s.db.Get(result, q, accountID, storageType); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbAccountStorageService) List(accountID uint32, mask uint8) ([]*AccountStorage, error) {
	var result []*AccountStorage
	var types []AccountStorageType

	if mask&uint8(AccountData) > 0 {
		types = append(types, AccountData)
	}
	if mask&uint8(AccountKeybinds) > 0 {
		types = append(types, AccountKeybinds)
	}
	if mask&uint8(AccountMacros) > 0 {
		types = append(types, AccountMacros)
	}

	if len(types) == 0 {
		return nil, nil
	}

	q, args, err := sqlx.In(`SELECT * FROM account_storage WHERE type IN (?) AND account_id = ?`, types, accountID)
	if err != nil {
		return nil, err
	}

	q = s.db.Rebind(q)
	if err := s.db.Select(&result, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (s *DbAccountStorageService) create(db creater, storage *AccountStorage) error {
	q := `
	INSERT INTO account_storage (account_id, type, data, uncompressed_size)
	VALUES (:account_id, :type, :data, :uncompressed_size)
	RETURNING updated_at`
	result, err := db.NamedQuery(q, storage)
	if err != nil {
		return err
	}

	result.Next()
	result.Scan(&storage.UpdatedAt)
	return nil
}

func (s *DbAccountStorageService) update(db updater, storage *AccountStorage) (bool, error) {
	q := `
	UPDATE account_storage
	SET data=:data, uncompressed_size=:uncompressed_size, updated_at=now()
	WHERE account_id=:account_id AND type=:type
	RETURNING updated_at`
	result, err := db.NamedQuery(q, storage)
	if err != nil {
		return false, err
	}

	// Was a row updated?
	if result.Next() {
		result.Scan(&storage.UpdatedAt)
		return true, nil
	}

	return false, nil
}

func (s *DbAccountStorageService) UpdateOrCreate(storage *AccountStorage) (bool, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return false, err
	}
	defer func() {
		tx.Rollback()
	}()

	updated, err := s.update(tx, storage)
	if err != nil {
		return false, err
	} else if !updated {
		if err := s.create(tx, storage); err != nil {
			return false, err
		}
	}

	return !updated, tx.Commit()
}

type CharacterStorageService interface {
	// Get returns the character storage with the given type. If it doesn't exist, the returned storage is nil.
	Get(uint32, CharacterStorageType) (*CharacterStorage, error)

	// List returns all storage belonging to characterID whose type matches mask. List may return fewer
	// results than specified in the mask. The results are not ordered.
	List(uint32, uint8) ([]*CharacterStorage, error)

	// UpdateOrCreate tries to update the storage if it exists, otherwise it's created. UpdateOrCreate
	// reports whether the storage was created. If the error is nil and created is false, then it
	// was updated.
	UpdateOrCreate(*CharacterStorage) (bool, error)
}

type DbCharacterStorageService struct {
	db *sqlx.DB
}

var _ CharacterStorageService = (*DbCharacterStorageService)(nil)

func NewDbCharacterStorageService(db *sqlx.DB) CharacterStorageService {
	return &DbCharacterStorageService{db}
}

func (s *DbCharacterStorageService) Get(characterID uint32, storageType CharacterStorageType) (*CharacterStorage, error) {
	result := &CharacterStorage{}
	q := `SELECT * FROM character_storage WHERE character_id = $1 AND type = $2`
	if err := s.db.Get(result, q, characterID, storageType); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbCharacterStorageService) List(characterID uint32, mask uint8) ([]*CharacterStorage, error) {
	var result []*CharacterStorage
	var types []CharacterStorageType

	if mask&uint8(CharacterConfig) > 0 {
		types = append(types, CharacterConfig)
	}
	if mask&uint8(CharacterKeybinds) > 0 {
		types = append(types, CharacterKeybinds)
	}
	if mask&uint8(CharacterMacros) > 0 {
		types = append(types, CharacterMacros)
	}
	if mask&uint8(CharacterLayout) > 0 {
		types = append(types, CharacterLayout)
	}
	if mask&uint8(CharacterChat) > 0 {
		types = append(types, CharacterChat)
	}

	if len(types) == 0 {
		return nil, nil
	}

	q, args, err := sqlx.In(`SELECT * FROM character_storage WHERE type IN (?) AND character_id = ?`, types, characterID)
	if err != nil {
		return nil, err
	}

	q = s.db.Rebind(q)
	if err := s.db.Select(&result, q, args...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}

func (s *DbCharacterStorageService) create(db creater, storage *CharacterStorage) error {
	q := `
	INSERT INTO character_storage (character_id, type, data, uncompressed_size)
	VALUES (:character_id, :type, :data, :uncompressed_size)
	RETURNING updated_at`
	result, err := db.NamedQuery(q, storage)
	if err != nil {
		return err
	}

	result.Next()
	result.Scan(&storage.UpdatedAt)
	return nil
}

func (s *DbCharacterStorageService) update(db updater, storage *CharacterStorage) (bool, error) {
	q := `
	UPDATE character_storage
	SET data=:data, uncompressed_size=:uncompressed_size, updated_at=now()
	WHERE character_id=:character_id AND type=:type
	RETURNING updated_at`
	result, err := db.NamedQuery(q, storage)
	if err != nil {
		return false, err
	}

	// Was a row updated?
	if result.Next() {
		result.Scan(&storage.UpdatedAt)
		return true, nil
	}

	return false, nil
}

func (s *DbCharacterStorageService) UpdateOrCreate(storage *CharacterStorage) (bool, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return false, err
	}
	defer func() {
		tx.Rollback()
	}()

	updated, err := s.update(tx, storage)
	if err != nil {
		return false, err
	} else if !updated {
		if err := s.create(tx, storage); err != nil {
			return false, err
		}
	}

	return !updated, tx.Commit()
}
