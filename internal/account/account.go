package account

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Account struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password"`
	ModifiedDate time.Time `json:"modified_date" db:"modified_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type AccountCredentialSpec struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Account) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}