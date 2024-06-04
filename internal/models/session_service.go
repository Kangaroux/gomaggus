package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type SessionService interface {
	// Get returns the matching Session, or nil if it doesn't exist.
	Get(uint32) (*Session, error)

	// Create creates a new Session and sets the Id and CreatedAt fields.
	Create(*Session) error

	// Update tries to update an existing Session and returns if it was updated.
	Update(*Session) (bool, error)

	// Delete tries to delete an existing Session by id and returns if it was deleted.
	Delete(uint32) (bool, error)
}

type DbSessionService struct {
	db *sqlx.DB
}

var _ SessionService = (*DbSessionService)(nil)

func NewDbSessionService(db *sqlx.DB) *DbSessionService {
	return &DbSessionService{db}
}

func (s *DbSessionService) Get(accountId uint32) (*Session, error) {
	result := &Session{}
	if err := s.db.Get(result, `SELECT * FROM sessions WHERE account_id = $1`, accountId); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}

func (s *DbSessionService) Create(session *Session) error {
	q := `
	INSERT INTO sessions (account_id, session_key, connected, connected_at, disconnected_at)
	VALUES (:account_id, :session_key, :connected, :connected_at, :disconnected_at)
	RETURNING id`
	result, err := s.db.NamedQuery(q, session)
	if err != nil {
		return err
	}
	result.Next()
	return result.StructScan(session)
}

func (s *DbSessionService) Update(session *Session) (bool, error) {
	q := `
	UPDATE sessions SET
	session_key=:session_key, connected=:connected, connected_at=:connected_at, disconnected_at=:disconnected_at
	WHERE id=:id`
	result, err := s.db.NamedExec(q, session)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}

func (s *DbSessionService) Delete(accountId uint32) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM sessions WHERE account_id=$1`, accountId)
	if err != nil {
		return false, err
	}
	n, _ := result.RowsAffected()
	return n > 0, err
}
