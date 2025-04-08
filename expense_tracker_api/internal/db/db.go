package wallet

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func InitDB(path string) *sqlx.DB {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	dir := wd + "/storage/"
	db, err := sqlx.Open("sqlite", dir+path)
	if err != nil {
		log.Fatal(err)
		return db
	}
	return db
}
