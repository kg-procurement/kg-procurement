package account

import "time"

type Account struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password"`
	ModifiedDate time.Time `json:"modified_date" db:"modified_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
