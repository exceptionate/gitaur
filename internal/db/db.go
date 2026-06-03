package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DBPath = "gitaur.db"

var Conn *sql.DB

func Init() error {
	var err error

	Conn, err = sql.Open("sqlite3", DBPath)
	if err != nil {
		return err
	}

	return nil
}

func Exists() bool {
	_, err := os.Stat(DBPath)
	return err == nil
}

func CreateSchema() error {
	_, err := Conn.Exec(UserSchema)
	if err != nil {
		return err
	}
	_, err = Conn.Exec(ProjectsSchema)
	if err != nil {
		return err
	}
	return err
}

func UserExists() (bool, error) {
	var count int

	err := Conn.QueryRow(
		"SELECT COUNT(*) FROM user",
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func SchemaExists() (bool, error) {
	var count int

	err := Conn.QueryRow(`
		SELECT COUNT(*)
		FROM sqlite_master
		WHERE type='table'
		AND name IN ('user', 'projects')
	`).Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 2, nil
}
