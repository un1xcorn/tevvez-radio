package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
}

func GetDatabase(name string) (Database, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return Database{}, err
	}

	return Database{db}, nil
}

func (db *Database) GetTracks() ([]string, error) {
	tracks := []string{}

	query, err := db.Query("SELECT * FROM tracks")
	if err != nil {
		return []string{}, err
	}

	for query.Next() {
		track := ""
		err = query.Scan(&track)
		if err != nil {
			return []string{}, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}
