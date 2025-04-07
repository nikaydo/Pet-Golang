package services

import (
	"database/sql"
	"errors"
	"main/internal/models"
	"main/internal/repositories"
	"strings"
)

var (
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type DBServ struct {
	Repo *repositories.Database
}

func NewUserService(repo *repositories.Database) *DBServ {
	return &DBServ{Repo: repo}
}

func (s *DBServ) GetUser(name string) (models.User, error) {
	u, err := s.Repo.GetUserByUsername(name)
	if err != nil {
		return u, err
	}
	err = u.DecryptRefresh()
	if err != nil {
		return u, err
	}
	return u, err
}

func (s *DBServ) AddUser(u models.User, b models.Balance) error {
	err := u.EncryptRefresh()
	if err != nil {
		return err
	}
	err = u.HashingPass()
	if err != nil {
		return err
	}
	res, err := s.Repo.AddUser(u)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrUserAlreadyExists
		}
		return err
	}
	u.ID = -1
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	u.ID = int(id)
	err = s.Repo.SetBalance(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBServ) IsUserExists(a models.Auth) (bool, models.User, error) {
	u, pass, err := s.Repo.UserExists(a)
	if err == sql.ErrNoRows {
		return false, u, nil
	}
	err = u.CheckPass(pass)
	if err != nil {
		return false, u, ErrInvalidPassword
	}
	err = u.DecryptRefresh()
	if err != nil {
		return false, u, err
	}
	return true, u, err
}

func (s *DBServ) UpdateBalance(b models.Balance, t string) error {
	return s.Repo.UpdateBalance(b, t)
}

func (s *DBServ) NewTransactions(t models.Transaction) error {
	b := models.Balance{
		UserID: t.UserID,
		Amount: t.Amount,
	}
	err := s.Repo.NewTransactions(t)
	if err != nil {
		return err
	}
	err = s.UpdateBalance(b, t.Type)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBServ) UpdateRefreshToken(u models.User) error {
	err := u.EncryptRefresh()
	if err != nil {
		return err
	}
	return s.Repo.UpdateRefreshToken(u)
}

func (s *DBServ) Transactions(id int) (models.Tlist, error) {
	return s.Repo.Transactions(id)
}

func (s *DBServ) Balance(id int) (float64, error) {
	return s.Repo.Balance(id)
}

func (s *DBServ) DelTrans(user_id, id int) error {
	return s.Repo.DelTrans(user_id, id)
}

func (s *DBServ) SearchTags(id int, tags ...string) (models.Tlist, error) {
	return s.Repo.SearchTags(id, tags...)
}
