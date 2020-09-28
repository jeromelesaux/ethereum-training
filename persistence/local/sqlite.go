package local

import "database/sql"

func ConnectSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "ethereum.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}
