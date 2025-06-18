package models

import (
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *User) ValidatePassword(plainText string) (bool, error) {
	// $2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj6kVKj1nHKS
	//err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	//if err != nil {
	//	switch {
	//	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
	//		return false, nil
	//	default:
	//		log.Println(err)
	//		return false, err
	//	}
	//}
	if u.Password != plainText {
		return false, errors.New("wrong password")
	}

	return true, nil
}
