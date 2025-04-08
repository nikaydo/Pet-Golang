package repositories

import (
	"errors"
	env "main/internal/config"
	db "main/internal/db"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Database struct {
	DB *sqlx.DB
}

func NewRepository(e *env.Config) *Database {
	return &Database{DB: db.InitDB(e.DB.Path)}
}

func (f *Database) MakeTable() error {
	_, err := f.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL, refresh_token TEXT NOT NULL);")
	if err != nil {
		return err
	}
	_, err = f.DB.Exec("CREATE TABLE IF NOT EXISTS balance (id INTEGER PRIMARY KEY AUTOINCREMENT, amount DECIMAL(10, 2) NOT NULL DEFAULT 0);")
	if err != nil {
		return err
	}
	_, err = f.DB.Exec("CREATE TABLE IF NOT EXISTS transactions (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, amount DECIMAL(10, 2) NOT NULL, type TEXT CHECK( type IN ('income', 'outcome') ) NOT NULL, date TEXT, note TEXT, tag TEXT);")
	if err != nil {
		return err
	}
	return nil
}

func (f *Database) Close() error {
	if f.DB != nil {
		return f.DB.Close()
	}
	return nil
}
