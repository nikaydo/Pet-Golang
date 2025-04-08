package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	env "main/internal/config"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidDecrypted = errors.New("invalid encrypted data")
)

func HashingPass(pass string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passHash), nil

}

func CheckPass(pass,passHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass))
	if err != nil {
		return err
	}
	return nil
}

func EncryptRefresh(refresh string, env env.Config) (string, error) {
	block, err := aes.NewCipher([]byte(env.Keys.SecretForAES))
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	cipherText := aesGCM.Seal(nil, nonce, []byte(refresh), nil)
	s := base64.StdEncoding.EncodeToString(append(nonce, cipherText...))
	return s, nil
}

func DecryptRefresh(refresh string, env env.Config) (string, error) {
	data, err := base64.StdEncoding.DecodeString(refresh)
	if err != nil {
		return "", err
	}
	if len(data) < 12+aes.BlockSize { // aes.BlockSize == 16
		return "", ErrInvalidDecrypted
	}
	nonce, cipherText := data[:12], data[12:]
	block, err := aes.NewCipher([]byte(env.Keys.SecretForAES))
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
