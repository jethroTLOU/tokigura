package database

import (
	"database/sql"
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func IsUniqueViolation(err error) bool {
	if sqliteErr, ok := errors.AsType[*sqlite.Error](err); ok {
		code := sqliteErr.Code()
		return code == sqlite3.SQLITE_CONSTRAINT_UNIQUE ||
			code == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
	}
	return false
}
