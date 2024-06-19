package database

import (
	"log"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)


var Db *sql.DB

const DbPath = "./pkg/database/database.db"

func init() {
	Db = databaseInit()
	//Migrate()
}

func databaseInit() *sql.DB {
	var err error
	Db, err = sql.Open("sqlite3", DbPath)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    log.Println("Database running")

	return Db
}