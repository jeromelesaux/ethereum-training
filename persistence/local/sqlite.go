package local

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "ethereum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
