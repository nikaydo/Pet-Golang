package models

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Balance struct {
	ID     int     `json:"-"`
	UserID int     `json:"-"`
	Amount float32 `json:"Amount,omitempty"`
}

type Transaction struct {
	ID     int     `json:"id,omitempty"`
	UserID int     `json:"-"`
	Amount float32 `json:"ammount"`
	Type   string  `json:"type"`
	Date   string  `json:"date,omitempty"`
	Note   string  `json:"note,omitempty"`
	Tag    string  `json:"tag,omitempty"`
}

type Tlist struct {
	Income  []Transaction `json:"income,omitempty"`
	Outcome []Transaction `json:"outcome,omitempty"`
}

type User struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
	Refresh  string `json:"-"`
}
