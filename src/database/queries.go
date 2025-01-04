package database

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func (dbConn *DBConnection) GetProcessedPaths() ([]string, error) {
	rows, err := dbConn.DB.Query("SELECT path FROM post_dirs WHERE updated = TRUE")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paths []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}

	return paths, nil
}

func (dbConn *DBConnection) UpdateProcessed(path string) error {
	// select id by path
	var id int
	err := dbConn.DB.QueryRow("SELECT id FROM post_dirs WHERE path = ?", path).Scan(&id)
	if err != nil {
		log.Printf("DB Error: %s", err)
	}

	if id == 0 {
		// if record not found, insert
		_, err = dbConn.DB.Exec("INSERT INTO post_dirs (path, updated, created_at, updated_at) VALUES (?, TRUE, ?, ?)", path, time.Now(), time.Now())
		if err != nil {
			return nil
		}
	} else {
		// update processed
		_, err = dbConn.DB.Exec("UPDATE post_dirs SET updated = TRUE, updated_at = ? WHERE id = ?", time.Now(), id)
		if err != nil {
			return err
		}
	}
	return nil
}
