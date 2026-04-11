package model

type Payment struct {
	ID     string `json:"id" db:"id"`
	Amount int    `json:"amount" db:"amount" validate:"required,gt=0"`
}
