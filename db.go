package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

const (
	cCreateTableStmt = "CREATE TABLE IF NOT EXISTS attendees(id INTEGER PRIMARY KEY, name VARCHAR, age INT)"
)

func createTable(db *sql.DB) error {
	statement, err := db.Prepare(cCreateTableStmt)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(); err != nil {
		return err
	}

	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "priv/sqlite.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// create Attendees table
	if err := createTable(db); err != nil {
		panic(err)
	}
	fmt.Println("Attendees table created")
}
