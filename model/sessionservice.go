package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SessionService interface {
	// Get returns a session by account id, or nil if it doesn't exist.
	Get(uint32) (*Session, error)

	// Create creates a new session and sets the Id and CreatedAt fields.
	Create(*Session) error

	// Update tries to update an existing session and returns if it was updated.
	Update(*Session) (bool, error)

	// Delete tries to delete an existing session by account id and returns if it was deleted.
	Delete(uint32) (bool, error)

	// UpdateOrCreate tries to update an existing session with the given account id. If no session
	// for that account exists, it creates one.
	UpdateOrCreate(*Session) error
}

type DbSessionService struct {
	db *sqlx.DB
}

var _ SessionService = (*DbSessionService)(nil)

func NewDbSessionService(db *sqlx.DB) *DbSessionService {
	return &DbSessionService{db}
}

func (s *DbSessionService) Get(accountId uint32) (*Session, error) {
	return s.get(s.db, accountId)
}

type getter interface {
	Get(dest interface{}, query string, args ...interface{}) error
}

func (s *DbSessionService) get(db getter, accountId uint32) (*Session, error) {
	result := &Session{}
	if err := db.Get(result, `SELECT * FROM sessions WHERE account_id = $1`, accountId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbSessionService) Create(session *Session) error {
	return s.create(s.db, session)
}

func (s *DbSessionService) create(db creater, session *Session) error {
	q := `
	INSERT INTO sessions (account_id, session_key, connected, connected_at, disconnected_at)
	VALUES (:account_id, :session_key, :connected, :connected_at, :disconnected_at)`
	_, err := db.NamedQuery(q, session)
	return err
}

func (s *DbSessionService) Update(session *Session) (bool, error) {
	return s.update(s.db, session)
}

func (s *DbSessionService) update(db updater, session *Session) (bool, error) {
	q := `
	UPDATE sessions
	SET session_key=:session_key, connected=:connected, connected_at=:connected_at, disconnected_at=:disconnected_at
	WHERE account_id=:account_id`

	result, err := db.NamedExec(q, session)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()

	return n > 0, err
}

func (s *DbSessionService) Delete(accountId uint32) (bool, error) {
	return s.delete(s.db, accountId)
}

type deleter interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func (s *DbSessionService) delete(db deleter, accountId uint32) (bool, error) {
	result, err := db.Exec(`DELETE FROM sessions WHERE account_id=$1`, accountId)
	if err != nil {
		return false, err
	}

	n, _ := result.RowsAffected()

	return n > 0, err
}

func (s *DbSessionService) UpdateOrCreate(session *Session) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	updated, err := s.update(tx, session)
	if err != nil {
		return err
	} else if !updated {
		if err := s.create(tx, session); err != nil {
			return err
		}
	}

	return tx.Commit()
}
