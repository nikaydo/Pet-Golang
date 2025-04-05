package wallet

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(path string) *sql.DB {
	dir := "../storage/"
	db, err := sql.Open("sqlite", dir+path)
	if err != nil {
		log.Fatal(err)
		return db
	}
	return db
}
