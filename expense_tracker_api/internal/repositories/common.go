package repositories

import (
	"database/sql"
	env "main/internal/config"
	"main/internal/models"
)

func getAESKey() ([]byte, error) {
	env := env.GetConfig()
	return []byte(env.Keys.SecretForAES), nil
}

func transForm(rows *sql.Rows) (models.Tlist, error) {
	var res models.Tlist
	for rows.Next() {
		t := models.Transaction{}
		err := rows.Scan(&t.ID, &t.Amount, &t.Type, &t.Date, &t.Note, &t.Tag)
		if err != nil {
			return models.Tlist{}, err
		}
		switch t.Type {
		case "outcome":
			res.Outcome = append(res.Outcome, t)
		case "income":
			res.Income = append(res.Income, t)
		}
	}
	return res, nil
}
