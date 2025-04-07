package models

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//id            int
	//RefreshToken string `json:"refresh_token"`
	//JwtToken      string `json:"jwt_token"`
	//role          string `json:"role"`
}

type Balance struct {
	ID     int     `json:"-"`
	UserID int     `json:"-"`
	Amount float32 `json:"Amount,omitempty"`
}

type Transaction struct {
	ID     int     `json:"-"`
	UserID int     `json:"-"`
	Amount float32 `json:"ammount"`
	Type   string  `json:"type"`
	Date   string  `json:"-"`
	Note   string  `json:"note,omitempty"`
	Tag    string  `json:"tag,omitempty"`
}

type Tlist struct {
	Income  []Transaction `json:"income,omitempty"`
	Outcome []Transaction `json:"outcome,omitempty"`
}
