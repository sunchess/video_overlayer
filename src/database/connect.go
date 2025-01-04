package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DBConnection struct {
	DBPath string
	DB     *sql.DB
}

func NewDBConnection(dbPath string) (*DBConnection, error) {
	db, err := initDB(dbPath)
	if err != nil {
		return nil, err
	}

	return &DBConnection{
		DBPath: dbPath,
		DB:     db,
	}, nil
}

func (dbConn *DBConnection) Close() error {
	return dbConn.DB.Close()
}

func initDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS post_dirs (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        path TEXT NOT NULL UNIQUE,
				processed BOOLEAN NOT NULL DEFAULT FALSE,
				updated Boolean NOT NULL DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = db.Exec("SELECT 1 FROM post_dirs LIMIT 1")
	if err != nil {
		if _, err := db.Exec(createTableSQL); err != nil {
			return nil, err
		}
	}
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}
