package repositories

import (
	env "main/internal/config"
)

func getAESKey() ([]byte, error) {
	env := env.GetConfig()
	return []byte(env.Keys.SecretForAES), nil
}
