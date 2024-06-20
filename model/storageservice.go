package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type AccountStorageService interface {
	// List returns all storage belonging to accountID whose type matches mask. List may return fewer
	// results than specified in the mask. The results are not ordered.
	List(uint32, uint8) ([]*AccountStorage, error)

	// UpdateOrCreate tries to update tje storage if it exists, otherwise it's created. UpdateOrCreate
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
	INSERT INTO account_storage (account_id, type, data)
	VALUES (:account_id, :type, :data)
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
	SET data=:data, updated_at=now()
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
