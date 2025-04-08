package jwt

import (
	"errors"
	"log"
	env "main/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired       = errors.New("token is expired")
	ErrExpNotFound        = errors.New("exp claim not found in token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrValidSigningMethod = errors.New("no valid signing method")
)

type NewJwtClaims struct {
	claims jwt.MapClaims
}

func (n *NewJwtClaims) setVar(id int, username string, role string, exTime time.Duration) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":      id,
			"username": username,
			"iss":      "server",
			"role":     role,
			"aud":      "money manager",
			"exp":      time.Now().Add(exTime).Unix(),
			"iat":      time.Now().Unix(),
		})
}

type JwtTokens struct {
	AccessToken   string
	AccessClaims  jwt.MapClaims
	RefreshToken  string
	RefreshClaims jwt.MapClaims
	Env           env.Config
}

func (j *JwtTokens) CreateTokens(id int, username string, role string) error {
	err := j.CreateJwtToken(id, username, role)
	if err != nil {
		log.Println("Error creating JWT token:", err)
		return err
	}
	err = j.CreateRefreshToken(id, username, role)
	if err != nil {
		log.Println("Error creating refresh token:", err)
		return err
	}

	return nil

}

func (j *JwtTokens) CreateJwtToken(id int, username string, role string) error {
	var claims NewJwtClaims
	cl := claims.setVar(id, username, role, j.Env.TTL.AccessToken)
	tokenString, err := cl.SignedString([]byte(j.Env.Keys.SecretKeyForJWT))
	if err != nil {
		return err
	}
	j.AccessToken = tokenString
	return nil
}

func (j *JwtTokens) CreateRefreshToken(id int, username string, role string) error {
	var claims NewJwtClaims
	cl := claims.setVar(id, username, role, j.Env.TTL.RefreshToken)
	tokenString, err := cl.SignedString([]byte(j.Env.Keys.SecretKeyForRefresh))
	if err != nil {
		return err
	}
	j.RefreshToken = tokenString
	return nil
}

func (j *JwtTokens) ValidateJwt() error {
	token, err := jwt.Parse(j.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrValidSigningMethod
		}
		return []byte(j.Env.Keys.SecretKeyForJWT), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			j.AccessClaims, err = setClaims(token)
			if err != nil {
				return err
			}
			return ErrTokenExpired
		}
		return err
	}
	j.AccessClaims, err = setClaims(token)
	if err != nil {
		return err
	}
	return nil
}

func (j *JwtTokens) ValidateRefresh() error {
	token, err := jwt.Parse(j.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrValidSigningMethod
		}
		return []byte(j.Env.Keys.SecretKeyForRefresh), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			j.RefreshClaims, err = setClaims(token)
			if err != nil {
				return err
			}
			return ErrTokenExpired
		}
		return err
	}
	j.RefreshClaims, err = setClaims(token)
	if err != nil {
		return err
	}
	return nil
}

func setClaims(token *jwt.Token) (jwt.MapClaims, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
