package model

type Payment struct {
	ID       string `json:"id" db:"id"`
	Amount   int    `json:"amount" db:"amount" validate:"required,gt=0"`
	Currency string `json:"currency" db:"currency" validate:"omitempty,currency"`
}
