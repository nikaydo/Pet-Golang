package services

import (
	j "main/internal/jwt"
)

type AuthServ struct {
	JwtClaims j.NewJwtClaims
	Jwt       j.JwtTokens
}

func NewAuthServ() *AuthServ {
	return &AuthServ{JwtClaims: j.NewJwtClaims{}, Jwt: j.JwtTokens{}}
}

func (s *AuthServ) CreateTokens(id int, username string, role string) error {
	return s.Jwt.CreateTokens(id, username, role)
}

func (s *AuthServ) CreateJwtToken(id int, username string, role string) error {
	return s.Jwt.CreateJwtToken(id, username, role)
}

func (s *AuthServ) CreateRefreshToken(id int, username string, role string) error {
	return s.Jwt.CreateRefreshToken(id, username, role)
}

func (s *AuthServ) ValidateJwt() error {
	return s.Jwt.ValidateJwt()
}

func (s *AuthServ) ValidateRefresh() error {
	return s.Jwt.ValidateRefresh()
}
