package services

import "main/internal/repositories"

type Handler struct {
	DBServ     *DBServ
	CryptoServ *CryptoServ
	AuthServ   *AuthServ
}

func NewHandler(repo *repositories.Database) *Handler {
	return &Handler{
		DBServ:     NewDBServ(repo),
		CryptoServ: NewCryptoServ(),
		AuthServ:   NewAuthServ(),
	}
}
