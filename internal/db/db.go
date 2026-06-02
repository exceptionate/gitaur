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
