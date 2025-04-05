package repositories

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Refresh  string `json:"-"`
}

func (u *User) HashingPass() error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(passHash)
	return nil

}

func (u *User) CheckPass(passHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(u.Password))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) encryptRefresh() error {
	key, err := getAESKey()
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	cipherText := aesGCM.Seal(nil, nonce, []byte(u.Refresh), nil)
	u.Refresh = base64.StdEncoding.EncodeToString(append(nonce, cipherText...))
	return nil
}

func (u *User) decryptRefresh() error {
	key, err := getAESKey()
	if err != nil {
		return err
	}
	data, err := base64.StdEncoding.DecodeString(u.Refresh)
	if err != nil {
		return err
	}
	if len(data) < 12+aes.BlockSize { // aes.BlockSize == 16
		return ErrInvalidDecrypted
	}
	nonce, cipherText := data[:12], data[12:]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return err
	}
	u.Refresh = string(plainText)
	return nil
}
