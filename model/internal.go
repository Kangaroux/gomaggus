package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type creater interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

type updater interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}
