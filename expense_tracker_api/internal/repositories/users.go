package repositories

import (
	"database/sql"
	"main/internal/models"
)

func (f *Database) AddUser(u models.User) (sql.Result, error) {
	new_user, err := f.DB.Exec("INSERT INTO users (username, password_hash,refresh_token) VALUES (:username, :password_hash, :refresh_token);",
		sql.Named("username", u.Username),
		sql.Named("password_hash", u.Password),
		sql.Named("refresh_token", u.Refresh))
	if err != nil {
		return nil, err
	}
	return new_user, nil
}
func (f *Database) UpdateRefreshToken(u models.User) error {
	_, err := f.DB.Exec("UPDATE users SET refresh_token = :refresh_token WHERE id = :id;",
		sql.Named("id", u.ID),
		sql.Named("refresh_token", u.Refresh))
	if err != nil {
		return err
	}
	return nil
} 
func (f *Database) UserExists(a models.Auth) (models.User, string, error) {
	u := models.User{Username: a.Username, Password: a.Password}
	row := f.DB.QueryRow("SELECT id, username,password_hash,refresh_token FROM users WHERE username = :username;",
		sql.Named("username", u.Username))
	var pass string
	err := row.Scan(&u.ID, &u.Username, &pass, &u.Refresh)
	if err != nil {
		return u, "", err
	}
	return u, pass, nil
}
func (f *Database) GetUserByUsername(username string) (models.User, error) {
	u := models.User{}
	row := f.DB.QueryRow("SELECT id, username,password_hash,refresh_token FROM users WHERE username = :username;",
		sql.Named("username", username))
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Refresh)
	if err != nil {
		return u, err
	}
	return u, nil
}
