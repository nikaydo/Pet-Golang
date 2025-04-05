package repositories

import (
	"database/sql"
	"errors"
	"main/internal/models"
	"strings"
	"time"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidDecrypted  = errors.New("invalid encrypted data")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
)

type File struct {
	DB *sql.DB
}

func NewRepository(DB *sql.DB) *File {
	return &File{DB: DB}
}

func (f *File) MakeTable() error {
	_, err := f.DB.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL, refresh_token TEXT NOT NULL);")
	if err != nil {
		return err
	}
	_, err = f.DB.Exec("CREATE TABLE IF NOT EXISTS balance (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, amount DECIMAL(10, 2) NOT NULL DEFAULT 0, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE);")
	if err != nil {
		return err
	}
	_, err = f.DB.Exec("CREATE TABLE IF NOT EXISTS transactions (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER NOT NULL, amount DECIMAL(10, 2) NOT NULL, type TEXT CHECK( type IN ('income', 'outcome') ) NOT NULL, date TEXT, note TEXT, tag TEXT, FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE );")
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Close() error {
	if f.DB != nil {
		return f.DB.Close()
	}
	return nil
}

func (f *File) setBalance(b models.Balance) error {
	_, err := f.DB.Exec("INSERT INTO balance (user_id, amount) VALUES (:user_id,:amount);",
		sql.Named("user_id", b.UserID),
		sql.Named("amount", b.Amount))
	if err != nil {
		return err
	}
	return nil
}

func (f *File) AddUser(u User, b models.Balance) error {
	err := u.encryptRefresh()
	if err != nil {
		return err
	}
	err = u.HashingPass()
	if err != nil {
		return err
	}
	u.ID = -1
	new_user, err := f.DB.Exec("INSERT INTO users (username, password_hash,refresh_token) VALUES (:username, :password_hash, :refresh_token);",
		sql.Named("username", u.Username),
		sql.Named("password_hash", u.Password),
		sql.Named("refresh_token", u.Refresh))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserAlreadyExists
		}
		return err
	}
	id, err := new_user.LastInsertId()
	if err != nil {
		panic(err)
	}
	u.ID = int(id)
	err = f.setBalance(b)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Balance(id int) (float64, error) {
	row := f.DB.QueryRow("SELECT amount FROM balance WHERE id = :user_id;", sql.Named("user_id", id))
	var ammount float64
	err := row.Scan(&ammount)

	if err != nil {
		return 0, err
	}
	return ammount, nil
}

func (f *File) NewTransactions(t models.Transaction) error {
	_, err := f.DB.Exec("INSERT INTO transactions (user_id, amount, type, note, date, tag) VALUES (:user_id, :amount, :type, :note, :date, :tag);",
		sql.Named("user_id", t.UserID),
		sql.Named("amount", t.Amount),
		sql.Named("date", time.Now().Format("02.01.2006 15:04:05")),
		sql.Named("note", t.Note),
		sql.Named("tag", t.Tag),
		sql.Named("type", t.Type))
	if err != nil {
		return err
	}
	b := models.Balance{
		UserID: t.UserID,
		Amount: t.Amount,
	}
	err = f.UpdateBalance(b, t.Type)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) UpdateBalance(b models.Balance, t string) error {
	var operation string
	switch t {
	case "outcome":
		operation = "-"
	case "income":
		operation = "+"
	}
	_, err := f.DB.Exec("UPDATE balance SET amount = amount "+operation+" :amount WHERE user_id = :user_id;",
		sql.Named("user_id", b.UserID),
		sql.Named("amount", b.Amount))
	if err != nil {
		return err
	}
	return nil
}

func (f *File) UpdateRefreshToken(u User) error {
	err := u.encryptRefresh()
	if err != nil {
		return err
	}
	_, err = f.DB.Exec("UPDATE users SET refresh_token = :refresh_token WHERE id = :id;",
		sql.Named("id", u.ID),
		sql.Named("refresh_token", u.Refresh))
	if err != nil {
		return err
	}
	return nil
}

func (f *File) Transactions(id int) (models.Tlist, error) {
	row, err := f.DB.Query("SELECT amount, type, date, note, tag FROM transactions WHERE user_id = :user_id;",
		sql.Named("user_id", id))
	if err != nil {
		return models.Tlist{}, err
	}
	var res models.Tlist
	for row.Next() {
		t := models.Transaction{}
		err := row.Scan(&t.Amount, &t.Type, &t.Date, &t.Note, &t.Tag)
		if err != nil {
			panic(err)
		}
		if t.Type == "outcome" {
			res.Outcome = append(res.Outcome, t)
		} else {
			res.Income = append(res.Income, t)
		}

	}
	return res, nil

}

func (f *File) IsUserExists(a models.Auth) (bool, User, error) {
	u := User{Username: a.Username, Password: a.Password}
	row := f.DB.QueryRow("SELECT id, username,password_hash,refresh_token FROM users WHERE username = :username;",
		sql.Named("username", u.Username))
	var pass string
	err := row.Scan(&u.ID, &u.Username, &pass, &u.Refresh)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, u, nil
		}
		return true, u, err
	}
	err = u.CheckPass(pass)
	if err != nil {
		return true, u, ErrInvalidPassword
	}
	err = u.decryptRefresh()
	if err != nil {
		return true, u, err
	}
	return true, u, nil
}

func (f *File) GetUserByUsername(username string) (User, error) {
	u := User{}
	row := f.DB.QueryRow("SELECT id, username,password_hash,refresh_token FROM users WHERE username = :username;",
		sql.Named("username", username))
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Refresh)
	if err != nil {
		return u, err
	}
	err = u.decryptRefresh()
	if err != nil {
		return u, err
	}
	return u, nil
}
