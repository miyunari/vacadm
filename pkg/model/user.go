package model

import "time"

type User struct {
	ID         string    `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Email      string    `json:"email"`
	CreatedAt  time.Time `json:"created_at"`
	DeletedAt  time.Time `json:"deleted_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (u User) Copy() *User {
	return &u
}
