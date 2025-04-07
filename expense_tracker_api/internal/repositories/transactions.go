package repositories

import (
	"database/sql"
	"main/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

func (f *Database) NewTransactions(t models.Transaction) error {
	_, err := f.DB.Exec("INSERT INTO transactions (user_id, amount, type, note, date, tag) VALUES (:user_id, :amount, :type, :note, :date, :tag);",
		sql.Named("user_id", t.UserID),
		sql.Named("amount", t.Amount),
		sql.Named("date", time.Now().Format("02.01.2006 15:04:05")),
		sql.Named("note", t.Note),
		sql.Named("tag", t.Tag),
		sql.Named("type", t.Type))
	if err != nil {
		return err
	}
	return nil
}
func (f *Database) Transactions(id int) (models.Tlist, error) {
	row, err := f.DB.Query("SELECT id, amount, type, date, note, tag FROM transactions WHERE user_id = :user_id;",
		sql.Named("user_id", id))
	if err != nil {
		return models.Tlist{}, err
	}
	var res models.Tlist
	for row.Next() {
		t := models.Transaction{}
		err := row.Scan(&t.ID, &t.Amount, &t.Type, &t.Date, &t.Note, &t.Tag)
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
func (f *Database) DelTrans(user_id, id int) error {
	_, err := f.DB.Exec("DELETE FROM transactions WHERE user_id = :user_id AND id = :id;",
		sql.Named("user_id", id),
		sql.Named("id", id))
	if err != nil {
		return err
	}
	return nil
}

func (f *Database) SearchTags(id int, tags ...string) (models.Tlist, error) {
	query := `SELECT id, amount, type, date, note, tag FROM transactions WHERE user_id = :user_id AND tag IN (:tags)`
	query, args, err := sqlx.Named(query, map[string]interface{}{
		"user_id": id,
		"tags":    tags,
	})
	if err != nil {
		return models.Tlist{}, err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return models.Tlist{}, err
	}
	query = sqlx.Rebind(sqlx.NAMED, query)
	rows, err := f.DB.Query(query, args...)
	if err != nil {
		return models.Tlist{}, err
	}
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
