package services

import (
	env "main/internal/config"
	n "main/internal/models"
)

type CryptoServ struct {
}

func NewCryptoServ() *CryptoServ {
	return &CryptoServ{}
}

func (s *CryptoServ) VerifykPass(pass, passHash string) error {
	return n.CheckPass(pass, passHash)
}

func (s *CryptoServ) HashingPass(pass string) (string, error) {
	return n.HashingPass(pass)
}

func (s *CryptoServ) EncryptRefresh(refresh string, env env.Config) (string, error) {
	return n.EncryptRefresh(refresh, env)
}

func (s *CryptoServ) DecryptRefresh(refresh string, env env.Config) (string, error) {
	return n.DecryptRefresh(refresh, env)
}
