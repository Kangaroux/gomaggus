package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type AccountStorageService interface {
	// List returns all storage belonging to accountID whose type matches mask. List may return fewer
	// results than specified in the mask. The results are not ordered.
	List(uint32, StorageMask) ([]*AccountStorage, error)

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
