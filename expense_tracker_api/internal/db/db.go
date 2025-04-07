package wallet

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func InitDB(path string) *sql.DB {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	dir := wd + "/storage/"
	db, err := sql.Open("sqlite", dir+path)
	if err != nil {
		log.Fatal(err)
		return db
	}
	return db
}
