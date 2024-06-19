package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccountStorageService interface {
	List(uint32, StorageMask) ([]*AccountStorage, error)
	UpdateOrCreate(*AccountStorage) error
}

type DbAccountStorageService struct {
	db *sqlx.DB
}

var _ AccountStorageService = (*DbAccountStorageService)(nil)

func NewDbAccountStorageService(db *sqlx.DB) *DbAccountStorageService {
	return &DbAccountStorageService{db}
}

func (s *DbAccountStorageService) List(accountID uint32, mask StorageMask) ([]*AccountStorage, error) {
	var result []*AccountStorage
	var types []StorageType

	// Convert type mask to discrete values
	for i := 0; i < 8; i++ {
		if mask%(1<<i) > 0 {
			types = append(types, StorageType(i))
		}
	}

	if len(types) == 0 {
		return nil, nil
	}

	q, args, err := sqlx.In(`SELECT * FROM account_storage WHERE type IN (?)`, types)
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
	q := `INSERT INTO account_storage (account_id, type, data) VALUES (:account_id, :type, :data)`
	if _, err := db.NamedQuery(q, storage); err != nil {
		return err
	}

	storage.UpdatedAt = time.Now()
	return nil
}

func (s *DbAccountStorageService) update(db updater, storage *AccountStorage) (bool, error) {
	q := `
	UPDATE account_storage
	SET data=:data, updated_at=now()
	WHERE account_id=:account_id AND type=:type`

	result, err := db.NamedExec(q, storage)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()
	if n > 0 {
		storage.UpdatedAt = time.Now()
	}

	return n > 0, nil
}

func (s *DbAccountStorageService) UpdateOrCreate(storage *AccountStorage) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	updated, err := s.update(tx, storage)
	if err != nil {
		return err
	} else if !updated {
		if err := s.create(tx, storage); err != nil {
			return err
		}
	}

	return tx.Commit()
}
