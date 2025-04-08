package repositories

import (
	"database/sql"
	"errors"
	"main/internal/models"
)

func (f *Database) SetBalance(b models.Balance) error {
	_, err := f.DB.Exec("INSERT INTO balance (amount) VALUES (:amount);",
		sql.Named("amount", b.Amount))
	if err != nil {
		return err
	}
	return nil
}
func (f *Database) Balance(id int) (float64, error) {
	row := f.DB.QueryRow("SELECT amount FROM balance WHERE id = :id;",
		sql.Named("id", id))
	var ammount float64
	err := row.Scan(&ammount)
	if err != nil {
		return 0, err
	}
	return ammount, nil
}
func (f *Database) UpdateBalance(b models.Balance, t string) error {
	var r string
	switch t {
	case "income":
		r = "UPDATE balance SET amount = amount + :amount WHERE id = :id;"
	case "outcome":
		r = "UPDATE balance SET amount = amount - :amount WHERE id = :id;"
	default:
		return errors.New("invalid operation")
	}
	_, err := f.DB.Exec(r,
		sql.Named("id", b.UserID),
		sql.Named("amount", b.Amount))
	if err != nil {
		return err
	}
	return nil
}
