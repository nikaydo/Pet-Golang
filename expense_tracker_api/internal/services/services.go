package services

import (
	"main/internal/models"
	"main/internal/repositories"
)

type UserService struct {
	Repo *repositories.File
}

func NewUserService(repo *repositories.File) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetUser(name string) (repositories.User, error) {
	return s.Repo.GetUserByUsername(name)
}

func (s *UserService) AddUser(u repositories.User, b models.Balance) error {
	return s.Repo.AddUser(u, b)
}

func (s *UserService) IsUserExists(a models.Auth) (bool, repositories.User, error) {
	return s.Repo.IsUserExists(a)
}
func (s *UserService) UpdateBalance(b models.Balance, t string) error {
	return s.Repo.UpdateBalance(b, t)
}

func (s *UserService) NewTransactions(t models.Transaction) error {
	return s.Repo.NewTransactions(t)
}

func (s *UserService) UpdateRefreshToken(u repositories.User) error {
	return s.Repo.UpdateRefreshToken(u)
}

func (s *UserService) Transactions(id int) (models.Tlist, error) {
	return s.Repo.Transactions(id)
}

func (s *UserService) Balance(id int) (float64, error) {
	return s.Repo.Balance(id)
}
